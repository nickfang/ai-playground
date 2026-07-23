import numpy as np
import potrace

mask = np.load("mask.npy").astype(bool)
h, w = mask.shape

bmp = potrace.Bitmap(~mask)
path = bmp.trace(turdsize=8, alphamax=1.0, opttolerance=0.2)

def f(v):
    s = f"{v:.2f}".rstrip("0").rstrip(".")
    return s if s else "0"

parts = []
for curve in path:
    p = curve.start_point; sx, sy = p.x, p.y
    parts.append(f"M{f(sx)} {f(sy)}")
    for seg in curve:
        p = seg.end_point; ex, ey = p.x, p.y
        if seg.is_corner:
            p = seg.c; cx, cy = p.x, p.y
            parts.append(f"L{f(cx)} {f(cy)}L{f(ex)} {f(ey)}")
        else:
            p = seg.c1; c1x, c1y = p.x, p.y
            p = seg.c2; c2x, c2y = p.x, p.y
            parts.append(f"C{f(c1x)} {f(c1y)} {f(c2x)} {f(c2y)} {f(ex)} {f(ey)}")
    parts.append("Z")

d = "".join(parts)
svg = (
    f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}" width="{w}" height="{h}">\n'
    f'<path d="{d}" fill="#111" fill-rule="evenodd" stroke="none"/>\n'
    f'</svg>\n'
)
with open("floorplan-walls.svg", "w") as fh:
    fh.write(svg)
print(f"curves: {sum(1 for _ in path)}  svg bytes: {len(svg)}")
