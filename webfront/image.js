var mosaics = [{
    "file": "json/goalie_2014.json",
    "name": "Goalie"
}, {
    "file": "json/goalie_cheer.json",
    "name": "Goalie Cheer"
}, {
    "file": "json/goal_cheer.json",
    "name": "Goal Cheer"
}, {
    "file": "json/snowman.json",
    "name": "Snowman"
}]

$(document).ready(function() {

    var mosaicLinkContainer = $("#mosaicLinks");

    $.each(mosaics, function(idx, mosaic) {
        mosaicLinkContainer.append('<li><a href="#" class="mosaicLink">' + mosaic.name + '</a></li>');
    });

    mosaicLinkContainer.on('click', '.mosaicLink', function(e) {
        var idx = $(this).parent().index();
        reset();
        loadMosaic(idx);
        e.preventDefault();
    });
    loadMosaic(1);


});

function reset() {
    //this is pretty hacky, but the quickest way to clear 
    //for the new image
    $('#canvas').remove();
    $('#imgArea').prepend('<canvas id="canvas"></canvas>');
    $("#tileImg").attr("src", "");
    $("#tileImgText").text("");
}

function loadMosaic(mosaicIndex) {
    $.get(mosaics[mosaicIndex].file, {})
        .done(function(data) {
            $("#canvas").width(data.Width);
            $("#canvas").height(data.Height);
            setupMosaic(data);
        });

}

function setupMosaic(photo) {
    var width = photo.Width;
    var height = photo.Height;
    var tileSize = photo.TileSize;
    var mosaicShowing = false;

    var numTiles = Math.ceil(height / tileSize) * Math.ceil(width / tileSize);

    var processedTiles = 0;

    var canvas = $("#canvas")[0];
    var context = canvas.getContext("2d");

    context.canvas.height = height;
    context.canvas.width = width;

    var loadedCount = 0;


    var topImage = new Image();
    var bottomImage = new Image();

    bottomImage.src = photo.mosaicImage;
    bottomImage.onload = pictureLoaded;
    topImage.src = photo.srcImage;
    topImage.onload = pictureLoaded;

    //canvas image fading/switching inspired by: http://jsfiddle.net/m1erickson/zw9S4/
    function pictureLoaded() {
        loadedCount++;
        if (loadedCount >= 2) {
            context.drawImage(topImage, 0, 0, width, height);
        }
    };


    $("#fade").click(function() {
        var tileCount = 0;
        for (var y = 0; y < height; y = y + tileSize) {
            for (var x = 0; x < width; x = x + tileSize) {
                setTimeout(startTileFade, 10 * tileCount++, x, y, 0);
            }
        }
    });

    $("#canvas").click(function(e) {
        if (!mosaicShowing) {
            return;
        }
        coords = getRelativeMousePos($(this)[0], e);
        var newX = Math.floor(coords.x / tileSize);
        var newY = Math.floor(coords.y / tileSize);
        var pos = newX + "," + newY;
        $("#tileImg").attr("src", photo.Tiles[pos].Photo.Url);
        $("#tileImgText").text(photo.Tiles[pos].Photo.Text);

    });

    function startTileFade(x, y, fadePct) {
        animateFade();

        function animateFade() {
            if (fadePct > 100) {
                tileProcessed();
                return;
            }
            requestAnimationFrame(animateFade);
            draw(bottomImage, x, y, fadePct / 100);
            draw(topImage, x, y, (1 - fadePct / 100));
            fadePct++;
        }
    }

    function tileProcessed() {
        processedTiles++;
        if (processedTiles >= numTiles) {

            processedTiles = 0;
            var temp = bottomImage;
            bottomImage = topImage;
            topImage = temp;

            mosaicShowing = !mosaicShowing;

        }

    }

    function draw(img, x, y, opacity) {
        context.save();
        context.globalAlpha = opacity;
        context.drawImage(img, x, y, tileSize, tileSize, x, y, tileSize, tileSize);
        context.restore();
    }

    function getRelativeMousePos(canvas, evt) {
        var rect = canvas.getBoundingClientRect();
        return {
            x: evt.clientX - rect.left,
            y: evt.clientY - rect.top
        };
    }

}