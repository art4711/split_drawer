## Split drawing ##

Splits operations in image/draw.Drawer into several goroutines for
images that are heavy for some reason (like generated on the fly).
