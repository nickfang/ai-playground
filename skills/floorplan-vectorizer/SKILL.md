---
name: floorplan-vectorizer
description: Convert a raster floorplan image (screenshot, scan, brochure page) into a clean vector SVG/PDF containing only the architecture — walls, stairs, columns, elevator shafts — with text labels, room color fills, and furniture removed. Use when the user wants to vectorize, trace, or clean up a floor plan image, or extract just the walls from one.
---

# Floorplan Vectorizer

Produce a pixel-faithful vector of a floorplan's architecture. Do NOT hand-trace
the image into SVG paths by eye — that yields a crude approximation. Extract the
real geometry from the pixels instead.

## Pipeline

Work in a scratch directory. Copy the scripts from `scripts/` next to the source
image, named `source.png`.

1. **Setup** (once per machine):
   ```sh
   python3 -m venv venv
   venv/bin/pip install pillow numpy opencv-python-headless potracer
   # plus rsvg-convert (brew install librsvg) for SVG→PDF/PNG
   ```

2. **Threshold + component analysis**: `venv/bin/python analyze.py` shows how
   the dark linework (`gray < 110`) splits into connected components. Colored
   fills and light backgrounds drop out at threshold time. Raise the threshold
   if the plan's lines are gray/faded.

3. **Read coordinates visually**: `venv/bin/python make_grid.py` writes
   `grid.png` with a labeled 100px grid — Read it to locate every text label and
   furniture item. For dense or rotated text, render zoomed crops:
   `make_grid.py x0 y0 x1 y1 [zoom]` → `crop.png` with a 20px grid.

4. **Clean**: edit the lists at the top of `extract.py`, then run it.
   - `TEXT` / `FURNITURE`: erase boxes `(x1,y1,x2,y2)`. A connected component is
     deleted only if its bbox is **fully contained** in a box, so walls crossing
     a box always survive — be generous. Rotated text dies to a normal box
     because each letter is its own small component.
   - `PIXEL_ERASE`: for letters/icons physically touching a wall (they merge
     into the wall's component, so containment can't catch them). Tight boxes.
   - `REPAIR`: `(x1,y1,x2,y2,thickness)` line segments to redraw wall pieces a
     pixel-erase clipped.

5. **Verify visually and iterate** — this is the heart of the skill. Read
   `removed.png` (kept = black, removed = red) and compare against the source:
   - Text/furniture still black → its bbox crossed a box edge (enlarge box) or
     it touches a wall (move to `PIXEL_ERASE` + `REPAIR`).
   - The script prints boxes with residual content; walls crossing a box are
     normal, merged letters also show up there.
   - Zoom into fixed regions with crops to confirm repairs; loop until clean.

6. **Vectorize**: `venv/bin/python vectorize.py` → `floorplan-walls.svg`, then:
   ```sh
   rsvg-convert -f pdf floorplan-walls.svg -o floorplan-walls.pdf
   rsvg-convert -b white -w 1600 floorplan-walls.svg -o preview.png
   ```
   Read `preview.png` side-by-side against the source before delivering.

## Keep vs remove

Keep (architecture, even though it can look like clutter): stair tread
hatching, structural columns (small solid squares/diamonds on walls), dashed
room-divider lines (tiny components — ensure no erase box fully contains the
dashes), elevator shafts with their "X" glyphs (standard notation, and fused to
the shaft outlines anyway).

Remove: all text, furniture (tables, chairs, consoles/credenzas — often
wall-attached, needing pixel erases), restroom person icons, decorative fills.

## Gotchas

- **potracer polarity**: pass the mask *inverted* (`~mask`) to
  `potrace.Bitmap()` — otherwise the background gets traced and you get a black
  canvas with white lines. `vectorize.py` already does this.
- potracer curve points are objects with `.x`/`.y`, not tuples.
- macOS screenshot filenames contain a narrow no-break space (U+202F) before
  "AM/PM" — a typed regular space won't match the path.
- The `extract.py` lists ship with values from the plan this was built on;
  they're a worked example to pattern-match, replace them wholesale.
