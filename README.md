# gtyper
[ttyper](https://github.com/max-niederman/ttyper) inspired clone


![gtyper_demo](https://github.com/user-attachments/assets/c9ada2ca-bd07-40af-8ac6-3e25df7a870b)

## Report

### Keypresses Report View
<img width="1748" height="672" alt="2026-03-25-233845_hyprshot" src="https://github.com/user-attachments/assets/ef706c84-aaa7-474a-bed5-ae647cacd716" />


### Time Series Report View
<img width="1746" height="667" alt="2026-03-25-233901_hyprshot" src="https://github.com/user-attachments/assets/130caa19-6e23-461b-85bc-b2ed054d5b22" />

At the end of each session you'll see:

- **Adjusted WPM** — Correct characters typed divided by 5, divided by elapsed minutes.
- **Raw WPM** — All keypresses (including errors) divided by 5, divided by elapsed minutes.
- **Accuracy** — Correct keypresses as a percentage of total keypresses.
- **Correct Keys** — Correct keypresses out of total keypresses (e.g. `272/284`).
- **Worst Keys** — Up to 5 characters with the lowest per-key accuracy.
- **Chart** — Rolling 10-keypress WPM plotted over the course of the session.

## Install

### Shell
#### curl
```
sh -c "$(curl -fsSL https://raw.githubusercontent.com/mikegetz/gtyper/main/tools/install.sh)"
```
#### wget
```
sh -c "$(wget -qO- https://raw.githubusercontent.com/mikegetz/gtyper/main/tools/install.sh)"
```

## Usage

After installing, launch from your terminal:

```
gtyper
```

Type the displayed prompt as accurately and quickly as possible. When finished, a results screen shows your stats. Press `esc` to quit.

## Prompts

Opening passages from classic novels:

- *A Tale of Two Cities* — Charles Dickens
- *Moby-Dick* — Herman Melville
- *Pride and Prejudice* — Jane Austen
- *Neuromancer* — William Gibson
- *The Hobbit* — J.R.R. Tolkien
- *1984* — George Orwell
- *The Great Gatsby* — F. Scott Fitzgerald
- *Anna Karenina* — Leo Tolstoy
- *The Old Man and the Sea* — Ernest Hemingway
- *East of Eden* — John Steinbeck
- *Ulysses* — James Joyce
- *To Kill a Mockingbird* — Harper Lee
- *The Sun Also Rises* — Ernest Hemingway
- *Brave New World* — Aldous Huxley
- *One Hundred Years of Solitude* — Gabriel Garcia Marquez
- *The Trial* — Franz Kafka
- *The Color Purple* — Alice Walker
