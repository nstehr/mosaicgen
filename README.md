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
 - github.com/lucasb-eyer/go-colorful

##Usage
- The project will produce two binaries, **mosaicgen** and **indexer**.  Run **indexer <keyword>** to build a collection of images (based on a keyword) to use as sources for the photomosaic.  Run **mosaicgen /path/to/source/img** to produce the photomosaic.

##TODO
- Try to make it faster
- Implement method to use a tile image only once

##Author

Nathan Stehr

##Thanks

I was able to work on a solid chunk of this project on CreativiDay time provided by my employer, Macadamian Technologies.
