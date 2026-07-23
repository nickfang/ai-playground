"""Overlay coordinate grids on source.png so erase-box coordinates can be read
off visually.

  python make_grid.py                  -> grid.png (full image, 100px grid)
  python make_grid.py x0 y0 x1 y1 [z]  -> crop.png (region at zoom z, default 3,
                                          20px grid) for tricky areas
"""
import sys
import cv2

src = cv2.imread("source.png")
if src is None:
    sys.exit("source.png not found")

if len(sys.argv) >= 5:
    x0, y0, x1, y1 = map(int, sys.argv[1:5])
    z = int(sys.argv[5]) if len(sys.argv) > 5 else 3
    img = src[y0:y1, x0:x1].copy()
    img = cv2.resize(img, None, fx=z, fy=z, interpolation=cv2.INTER_NEAREST)
    step = 20
    h, w = img.shape[:2]
    for x in range(0, w, step * z):
        cv2.line(img, (x, 0), (x, h), (255, 0, 0), 1)
        cv2.putText(img, str(x0 + x // z), (x + 2, 20),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.4, (255, 0, 0), 1)
    for y in range(0, h, step * z):
        cv2.line(img, (0, y), (w, y), (255, 0, 0), 1)
        cv2.putText(img, str(y0 + y // z), (2, y + 14),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.4, (255, 0, 0), 1)
    cv2.imwrite("crop.png", img)
    print(f"wrote crop.png ({x0},{y0})-({x1},{y1}) zoom {z}")
else:
    img = src.copy()
    h, w = img.shape[:2]
    for x in range(0, w, 100):
        cv2.line(img, (x, 0), (x, h), (255, 0, 0), 2)
        cv2.putText(img, str(x), (x + 4, 30),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.9, (255, 0, 0), 2)
    for y in range(0, h, 100):
        cv2.line(img, (0, y), (w, y), (255, 0, 0), 2)
        cv2.putText(img, str(y), (4, y + 28),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.9, (255, 0, 0), 2)
    cv2.imwrite("grid.png", img)
    print(f"wrote grid.png ({w}x{h})")
