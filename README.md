# greentext CLI

The Green-Text Meme Generator CLI is a command-line tool built with Go that lets you create classic green-text memes with ease.

## Features

- **Customizable Inputs**: Write your own green-text stories with ease.
- **Thumbnail Support**: Add an image to your meme, or use the built-in placeholder if you're feeling minimalist.
- **Fast and Lightweight**: Powered by Go, because performance matters even when making memes.
- **Formatted Output**: Saves memes as beautifully formatted PNGs, ready for sharing.

## Installation

1. Make sure you have [Go](https://golang.org/dl/) installed on your machine.
2. Clone the repository:

   ```bash
   git clone https://github.com/jasonuc/greentext.git
   ```

3. Navigate to the project directory:

   ```bash
   cd greentext
   ```

4. Build the CLI tool:

   ```bash
   go build -o bin/greentext .
   ```

5. Run the tool:

   ```bash
   bin/greentext -h
   ```

## Usage

Here’s how you can generate your very own green-text meme:

### Example Command

```bash
bin/greentext -l 5 -t ./tfw.png -o meme.png
```

- `-l`: Number of lines to include in the green-text.
- `-t`: Path to a thumbnail image (optional).
- `-o`: Output file name.

### Example Input

When prompted, enter your meme lines:

```text
> be me
> decide to write a meme generator in Go
> "how hard could image manipulation be?"
> 3 hours in, still trying to draw a rectangle
> find 5 different libraries that almost do what I want
> none of them support rounded corners
> consider switching to Python
> remember I already told everyone Go is the best
> silently suffer
> finally get it working
> "go build -o bin/greentext"
> runtime panic: freetype called with nil font
> tfw image manipulation in Go makes you question life
```

### Example Output

This command will generate a beautiful green-text meme like this:

![Example Meme](example-meme.png)

## Contributing

Contributions are welcome! Whether you find a bug, have an idea for a new feature, or just want to improve the documentation, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

Happy meme-ing! 🚀
