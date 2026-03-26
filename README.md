# gtyper
[ttyper](https://github.com/max-niederman/ttyper) inspired clone

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

## Report

At the end of each session you'll see:

- **Adjusted WPM** — Correct characters typed divided by 5, divided by elapsed minutes.
- **Raw WPM** — All keypresses (including errors) divided by 5, divided by elapsed minutes.
- **Accuracy** — Correct keypresses as a percentage of total keypresses.
- **Correct Keys** — Correct keypresses out of total keypresses (e.g. `272/284`).
- **Worst Keys** — Up to 5 characters with the lowest per-key accuracy.
- **Chart** — Rolling 10-keypress WPM plotted over the course of the session.

## Prompts

Opening passages from classic novels:

- *A Tale of Two Cities* — Charles Dickens
- *Moby-Dick* — Herman Melville
- *Pride and Prejudice* — Jane Austen
- *Neuromancer* — William Gibson
- *The Hobbit* — J.R.R. Tolkien
