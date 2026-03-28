package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type book struct {
	title  string
	author string
}

// Project Gutenberg book IDs mapped to title and author.
// Text fetched directly from gutenberg.org cache.
var books = map[int]book{
	// Jane Austen
	161:  {"Sense and Sensibility", "Jane Austen"},
	1342: {"Pride and Prejudice", "Jane Austen"},
	158:  {"Emma", "Jane Austen"},
	105:  {"Persuasion", "Jane Austen"},
	141:  {"Mansfield Park", "Jane Austen"},
	121:  {"Northanger Abbey", "Jane Austen"},

	// Charles Dickens
	98:   {"A Tale of Two Cities", "Charles Dickens"},
	1400: {"Great Expectations", "Charles Dickens"},
	766:  {"David Copperfield", "Charles Dickens"},
	730:  {"Oliver Twist", "Charles Dickens"},
	1023: {"Bleak House", "Charles Dickens"},
	883:  {"Our Mutual Friend", "Charles Dickens"},
	967:  {"Nicholas Nickleby", "Charles Dickens"},
	821:  {"Dombey and Son", "Charles Dickens"},
	700:  {"The Old Curiosity Shop", "Charles Dickens"},
	917:  {"Barnaby Rudge", "Charles Dickens"},
	968:  {"Martin Chuzzlewit", "Charles Dickens"},

	// Mark Twain
	74:   {"The Adventures of Tom Sawyer", "Mark Twain"},
	76:   {"Adventures of Huckleberry Finn", "Mark Twain"},
	86:   {"A Connecticut Yankee in King Arthur's Court", "Mark Twain"},
	245:  {"Life on the Mississippi", "Mark Twain"},
	3176: {"The Man That Corrupted Hadleyburg", "Mark Twain"},
	102:  {"The Prince and the Pauper", "Mark Twain"},

	// Arthur Conan Doyle
	244:  {"A Study in Scarlet", "Arthur Conan Doyle"},
	1661: {"The Adventures of Sherlock Holmes", "Arthur Conan Doyle"},
	2852: {"The Hound of the Baskervilles", "Arthur Conan Doyle"},
	108:  {"The Sign of the Four", "Arthur Conan Doyle"},
	2097: {"The Memoirs of Sherlock Holmes", "Arthur Conan Doyle"},
	834:  {"The Return of Sherlock Holmes", "Arthur Conan Doyle"},

	// H.G. Wells
	35:   {"The Time Machine", "H. G. Wells"},
	36:   {"The War of the Worlds", "H. G. Wells"},
	159:  {"The Island of Doctor Moreau", "H. G. Wells"},
	5230: {"The Invisible Man", "H. G. Wells"},
	1013: {"The First Men in the Moon", "H. G. Wells"},
	2912: {"Tono-Bungay", "H. G. Wells"},

	// Jules Verne
	103:  {"Around the World in Eighty Days", "Jules Verne"},
	164:  {"Twenty Thousand Leagues Under the Sea", "Jules Verne"},
	1268: {"The Mysterious Island", "Jules Verne"},
	3748: {"A Journey to the Centre of the Earth", "Jules Verne"},
	3526: {"From the Earth to the Moon", "Jules Verne"},

	// Gothic / horror
	345:   {"Dracula", "Bram Stoker"},
	84:    {"Frankenstein", "Mary Wollstonecraft Shelley"},
	43:    {"The Strange Case of Dr Jekyll and Mr Hyde", "Robert Louis Stevenson"},
	696:   {"The Castle of Otranto", "Horace Walpole"},
	2428:  {"The Mysteries of Udolpho", "Ann Radcliffe"},
	4380:  {"The Monk", "Matthew Lewis"},
	18247: {"The Last Man", "Mary Shelley"},

	// Robert Louis Stevenson
	120: {"Treasure Island", "Robert Louis Stevenson"},
	419: {"Kidnapped", "Robert Louis Stevenson"},
	392: {"The Black Arrow", "Robert Louis Stevenson"},

	// Oscar Wilde
	174: {"The Picture of Dorian Gray", "Oscar Wilde"},

	// Jack London
	215:  {"The Call of the Wild", "Jack London"},
	910:  {"White Fang", "Jack London"},
	1164: {"The Sea-Wolf", "Jack London"},

	// Joseph Conrad
	219:  {"Heart of Darkness", "Joseph Conrad"},
	974:  {"Lord Jim", "Joseph Conrad"},
	2021: {"The Portrait of a Lady", "Henry James"},

	// Thomas Hardy
	110: {"Tess of the d'Urbervilles", "Thomas Hardy"},
	107: {"Far from the Madding Crowd", "Thomas Hardy"},
	143: {"The Mayor of Casterbridge", "Thomas Hardy"},
	122: {"The Return of the Native", "Thomas Hardy"},
	153: {"Jude the Obscure", "Thomas Hardy"},

	// George Eliot
	145:  {"Middlemarch", "George Eliot"},
	550:  {"Silas Marner", "George Eliot"},
	507:  {"Adam Bede", "George Eliot"},
	6688: {"The Mill on the Floss", "George Eliot"},

	// Brontë sisters
	768:  {"Wuthering Heights", "Emily Brontë"},
	1260: {"Jane Eyre", "Charlotte Brontë"},
	969:  {"The Tenant of Wildfell Hall", "Anne Brontë"},

	// American classics
	514:   {"Little Women", "Louisa May Alcott"},
	526:   {"The Scarlet Letter", "Nathaniel Hawthorne"},
	205:   {"Walden", "Henry David Thoreau"},
	1952:  {"The Yellow Wallpaper", "Charlotte Perkins Gilman"},
	8800:  {"The Souls of Black Folk", "W. E. B. Du Bois"},
	7370:  {"Anthem", "Ayn Rand"},
	203:   {"Uncle Tom's Cabin", "Harriet Beecher Stowe"},
	2701:  {"Moby-Dick", "Herman Melville"},
	10712: {"Bartleby the Scrivener", "Herman Melville"},
	71:    {"The Red Badge of Courage", "Stephen Crane"},
	160:   {"The Awakening", "Kate Chopin"},
	15130: {"Herland", "Charlotte Perkins Gilman"},
	2148:  {"O Pioneers!", "Willa Cather"},
	2660:  {"My Ántonia", "Willa Cather"},
	17396: {"Main Street", "Sinclair Lewis"},
	1015:  {"Sister Carrie", "Theodore Dreiser"},
	23:    {"Narrative of the Life of Frederick Douglass", "Frederick Douglass"},
	2376:  {"Up from Slavery", "Booker T. Washington"},

	// Lewis Carroll
	11: {"Alice's Adventures in Wonderland", "Lewis Carroll"},
	12: {"Through the Looking-Glass", "Lewis Carroll"},

	// L.M. Montgomery
	45: {"Anne of Green Gables", "L. M. Montgomery"},
	47: {"Anne of Avonlea", "L. M. Montgomery"},
	51: {"Anne of the Island", "L. M. Montgomery"},

	// L. Frank Baum
	55:    {"The Wonderful Wizard of Oz", "L. Frank Baum"},
	54:    {"The Marvelous Land of Oz", "L. Frank Baum"},
	22566: {"Ozma of Oz", "L. Frank Baum"},

	// Rudyard Kipling
	140:  {"The Jungle Book", "Rudyard Kipling"},
	236:  {"The Second Jungle Book", "Rudyard Kipling"},
	2226: {"Kim", "Rudyard Kipling"},
	2778: {"Just So Stories", "Rudyard Kipling"},

	// E.M. Forster
	624:  {"A Room with a View", "E. M. Forster"},
	2891: {"Howards End", "E. M. Forster"},

	// Henry James
	432:  {"The Turn of the Screw", "Henry James"},
	7178: {"The Age of Innocence", "Edith Wharton"},

	// Edith Wharton
	4517: {"Ethan Frome", "Edith Wharton"},
	284:  {"The House of Mirth", "Edith Wharton"},

	// F. Scott Fitzgerald
	17192: {"The Beautiful and Damned", "F. Scott Fitzgerald"},
	22381: {"This Side of Paradise", "F. Scott Fitzgerald"},
	64317: {"The Great Gatsby", "F. Scott Fitzgerald"},

	// James Joyce
	2814: {"Dubliners", "James Joyce"},
	4300: {"Ulysses", "James Joyce"},

	// Fyodor Dostoevsky
	2554:  {"Crime and Punishment", "Fyodor Dostoevsky"},
	2638:  {"The Idiot", "Fyodor Dostoevsky"},
	28054: {"The Brothers Karamazov", "Fyodor Dostoevsky"},
	600:   {"Notes from Underground", "Fyodor Dostoevsky"},
	8116:  {"The Possessed", "Fyodor Dostoevsky"},

	// Leo Tolstoy
	2600: {"War and Peace", "Leo Tolstoy"},
	1399: {"Anna Karenina", "Leo Tolstoy"},
	243:  {"The Death of Ivan Ilyich", "Leo Tolstoy"},
	689:  {"The Kreutzer Sonata", "Leo Tolstoy"},

	// Alexandre Dumas
	863:  {"The Count of Monte Cristo", "Alexandre Dumas"},
	1257: {"The Three Musketeers", "Alexandre Dumas"},
	1258: {"Twenty Years After", "Alexandre Dumas"},

	// Victor Hugo
	135:  {"Les Misérables", "Victor Hugo"},
	2610: {"The Hunchback of Notre-Dame", "Victor Hugo"},

	// Gustave Flaubert
	2413: {"Madame Bovary", "Gustave Flaubert"},

	// Guy de Maupassant
	3090: {"Bel-Ami", "Guy de Maupassant"},

	// Ivan Turgenev
	977: {"Fathers and Sons", "Ivan Turgenev"},

	// Nikolai Gogol
	1197: {"Dead Souls", "Nikolai Gogol"},

	// Franz Kafka
	5200: {"The Metamorphosis", "Franz Kafka"},

	// Hermann Hesse
	8500: {"Siddhartha", "Hermann Hesse"},

	// Voltaire
	19942: {"Candide", "Voltaire"},

	// Miguel de Cervantes
	996: {"Don Quixote", "Miguel de Cervantes"},

	// Frances Hodgson Burnett
	113: {"The Secret Garden", "Frances Hodgson Burnett"},
	479: {"Little Lord Fauntleroy", "Frances Hodgson Burnett"},

	// Daniel Defoe
	521: {"Robinson Crusoe", "Daniel Defoe"},
	376: {"Moll Flanders", "Daniel Defoe"},

	// Jonathan Swift
	829: {"Gulliver's Travels", "Jonathan Swift"},

	// Elizabeth Gaskell
	4276: {"North and South", "Elizabeth Gaskell"},
	2273: {"Wives and Daughters", "Elizabeth Gaskell"},
	4533: {"Cranford", "Elizabeth Gaskell"},

	// Anthony Trollope
	3271: {"Barchester Towers", "Anthony Trollope"},
	5231: {"The Way We Live Now", "Anthony Trollope"},

	// William Makepeace Thackeray
	599: {"Vanity Fair", "William Makepeace Thackeray"},

	// Washington Irving
	2048:  {"The Legend of Sleepy Hollow", "Washington Irving"},
	13514: {"Rip Van Winkle", "Washington Irving"},

	// Edgar Allan Poe
	2147: {"The Works of Edgar Allan Poe, Vol. I", "Edgar Allan Poe"},
	932:  {"The Fall of the House of Usher", "Edgar Allan Poe"},

	// Nathaniel Hawthorne
	512:  {"The House of the Seven Gables", "Nathaniel Hawthorne"},
	9257: {"Twice-Told Tales", "Nathaniel Hawthorne"},

	// Wilkie Collins
	583:  {"The Woman in White", "Wilkie Collins"},
	1653: {"The Moonstone", "Wilkie Collins"},

	// Baroness Orczy
	25344: {"The Scarlet Pimpernel", "Baroness Orczy"},

	// John Buchan
	558: {"The Thirty-Nine Steps", "John Buchan"},

	// H. Rider Haggard
	2166: {"King Solomon's Mines", "H. Rider Haggard"},
	3154: {"She", "H. Rider Haggard"},

	// P.G. Wodehouse
	7471:  {"My Man Jeeves", "P. G. Wodehouse"},
	10554: {"Something New", "P. G. Wodehouse"},

	// Jerome K. Jerome
	308: {"Three Men in a Boat", "Jerome K. Jerome"},

	// G.K. Chesterton
	1695: {"The Man Who Was Thursday", "G. K. Chesterton"},

	// Gaston Leroux
	175: {"The Phantom of the Opera", "Gaston Leroux"},

	// Maurice Leblanc
	6232: {"Arsène Lupin, Gentleman-Burglar", "Maurice Leblanc"},

	// D.H. Lawrence
	217:   {"Sons and Lovers", "D. H. Lawrence"},
	28948: {"Women in Love", "D. H. Lawrence"},

	// John Milton
	26: {"Paradise Lost", "John Milton"},

	// John Bunyan
	131: {"The Pilgrim's Progress", "John Bunyan"},

	// Oliver Goldsmith
	2587: {"The Vicar of Wakefield", "Oliver Goldsmith"},

	// Walter Scott
	5765: {"Ivanhoe", "Walter Scott"},

	// John Galsworthy
	2507: {"The Man of Property", "John Galsworthy"},

	// Giovanni Boccaccio
	3726: {"The Decameron", "Giovanni Boccaccio"},

	// Ambrose Bierce
	4362: {"An Occurrence at Owl Creek Bridge", "Ambrose Bierce"},

	// O. Henry
	2776: {"The Gift of the Magi and Other Stories", "O. Henry"},

	// Anthony Hope
	95: {"The Prisoner of Zenda", "Anthony Hope"},

	// William Godwin
	2402: {"Caleb Williams", "William Godwin"},

	// Upton Sinclair
	30254: {"The Jungle", "Upton Sinclair"},

	// Kate Chopin — 160 already added above

	// Stephen Crane
	8337: {"Maggie: A Girl of the Streets", "Stephen Crane"},

	// Ford Madox Ford
	2887: {"The Good Soldier", "Ford Madox Ford"},
}

type promptFetchedMsg struct {
	p   prompt
	err error
}

func fetchGutenbergPromptCmd() tea.Msg {
	ids := make([]int, 0, len(books))
	for id := range books {
		ids = append(ids, id)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	tried := map[int]bool{}
	for attempt := 0; attempt < 8; attempt++ {
		var id int
		for {
			id = ids[rand.Intn(len(ids))]
			if !tried[id] {
				break
			}
		}
		tried[id] = true

		book := books[id]
		url := fmt.Sprintf("https://www.gutenberg.org/cache/epub/%d/pg%d.txt", id, id)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil || resp.StatusCode != 200 {
			continue
		}

		p, err := extractPassage(body, book.title, book.author)
		if err != nil {
			continue
		}
		return promptFetchedMsg{p: p}
	}

	return promptFetchedMsg{err: errors.New("gutenberg: failed to fetch a suitable passage after 8 attempts")}
}

// normalizeText replaces Unicode punctuation that can't be typed on a standard
// keyboard with their closest ASCII equivalents.
func normalizeText(s string) string {
	// Gutenberg uses _word_ for italics and /# ... #/ for block annotations — strip them
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "/#", "")
	s = strings.ReplaceAll(s, "#/", "")

	r := strings.NewReplacer(
		"\u2018", "'", // left single quotation mark
		"\u2019", "'", // right single quotation mark (apostrophe)
		"\u201C", "\"", // left double quotation mark
		"\u201D", "\"", // right double quotation mark
		"\u2011", "-", // non-breaking hyphen
		"\u2013", "-", // en-dash
		"\u2014", "--", // em-dash
		"\u2026", "...", // horizontal ellipsis
		"\u00A0", " ", // non-breaking space
		"\u2012", "-", // figure dash
		"\u2015", "--", // horizontal bar
	)
	return r.Replace(s)
}

func punctuationDensity(s string) float64 {
	heavy := 0
	for _, ch := range s {
		switch ch {
		case '"', '!', '?', ';':
			heavy++
		}
	}
	return float64(heavy) / float64(len([]rune(s)))
}

func extractPassage(body []byte, title, author string) (prompt, error) {
	text := strings.ReplaceAll(string(body), "\r\n", "\n") // thanks windows
	text = normalizeText(text)

	const startMarker = "*** START OF THE PROJECT GUTENBERG EBOOK"
	if idx := strings.Index(text, startMarker); idx != -1 {
		nl := strings.Index(text[idx:], "\n")
		if nl != -1 {
			text = text[idx+nl+1:]
		}
	}

	const endMarker = "*** END OF"
	if idx := strings.Index(text, endMarker); idx != -1 {
		text = text[:idx]
	}

	paragraphs := strings.Split(text, "\n\n")
	var candidates []string
	for _, para := range paragraphs {
		normalized := strings.Join(strings.Fields(para), " ")
		n := len([]rune(normalized))
		if n < 200 || n > 500 {
			continue
		}
		if punctuationDensity(normalized) > 0.05 {
			continue
		}
		candidates = append(candidates, normalized)
	}

	if len(candidates) == 0 {
		return prompt{}, fmt.Errorf("gutenberg: no suitable passage found in %q", title)
	}

	citation := fmt.Sprintf("%s — %s", title, author)
	return prompt{text: candidates[rand.Intn(len(candidates))], source: citation}, nil
}
