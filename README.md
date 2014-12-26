#Photomosaic Generator

System to generate photomosaics (http://en.wikipedia.org/wiki/Photographic_mosaic).  A complete system that contains not only the image processing/generating functionality,
and image indexer to collect images from various public sources online  

##Requirements
 - Go (https://golang.org/)
 - Go resize (https://github.com/nfnt/resize)
 - https://bitbucket.org/liamstask/go-imgur
 - https://github.com/gedex/go-instagram
 - https://github.com/manki/flickgo
 - https://github.com/ChimeraCoder/anaconda
 - https://github.com/lucasb-eyer/go-colorful

##Usage
- The project will produce two binaries, **mosaicgen** and **indexer**.  Run **indexer <keyword>** to build a collection of images (based on a keyword) to use as sources for the photomosaic.  Run **mosaicgen /path/to/source/img** to produce the photomosaic.

##Examples
### Source Image
![Source Image](/examples/snowman1.png?raw=true "Source Image")
### Tiles Based On Average Colour
![Average Colour Tiler](/examples/tiled.png?raw=true "Average Colour Tiler")
### Tiles Based on algorithm by Matt Cutts (https://www.mattcutts.com/blog/photo-mosaic-effect-with-go/)
![Matt Cutts Inspired Tiler](/examples/tiled2.png?raw=true "Matt Cutts Inspired Tiler")
### Photomosaic (tile image can be reused)
![Photomosaic Tiler](/examples/tiled3_d.png?raw=true "Photomosaic Tiler")

##TODO
- Try to make it faster
- Implement method to use a tile image only once

##Author

Nathan Stehr

##Thanks

I was able to work on a solid chunk of this project on CreativiDay time provided by my employer, Macadamian Technologies.
