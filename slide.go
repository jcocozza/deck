package main

import "image"

/*
sample presentation file:

# 1
this is a presentation

# 2
It has content in it.
It also contains files!
@@file.png

# 3
// comments, and empty slides are allowed too

*/

type Slide struct {
	Title   string
	Content []string
	Image image.Image
}

var slides = []Slide{
	Slide{Title: "1", Content: []string{"foo bar baz"}},
	Slide{Title: "2", Content: []string{"foo", "bar a;sldkfja;lksjfdghj ;akljsch g;lkasjkdhnfg;lkjas  hgfv;lka js;lkja jsh;ljkash d;jkashd g;kasjdf ;lkasjdfg;l kasf;dlkg ja;sldkfgj a;slkdjfjgf ;askjdfjgf ;alskdjgf ;alskdfjg ;alskdfjfg ;laskdfjg", "baz", "", "MOO"}},
	Slide{Title: "3", Content: []string{"here is an image"}, Image: GenerateSampleImage()},
}
