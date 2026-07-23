"""Cleaning pass: threshold linework, delete text/furniture, keep architecture.

Edit the TEXT / FURNITURE / PIXEL_ERASE / REPAIR lists per plan (see README.md).
Current values are the worked example from the hotel conference-floor plan.
Outputs: filtered.png (kept), removed.png (kept=black/removed=red), mask.npy.
"""
import cv2
import numpy as np

src = cv2.imread("source.png")
gray = cv2.cvtColor(src, cv2.COLOR_BGR2GRAY)
mask = (gray < 110).astype(np.uint8)

# Erase boxes (x1,y1,x2,y2): remove connected components whose bbox lies FULLY inside.
# Walls crossing a box extend beyond it, so they are never removed.
TEXT = [
    (760, 85, 1000, 160),    # MANSFIELD A
    (645, 145, 850, 215),    # MANSFIELD
    (530, 255, 760, 325),    # MANSFIELD B
    (1165, 65, 1340, 150),   # SALES & CATERING
    (1030, 170, 1230, 255),  # EXECUTIVE OFFICES
    (295, 555, 555, 620),    # PENNYBACKER
    (110, 775, 275, 840),    # PFLUGER
    (700, 735, 815, 915),    # ELEVATORS (rotated)
    (340, 1115, 465, 1175),  # FOYER
    (660, 1095, 810, 1260),  # BALCONY (rotated)
    (110, 1175, 245, 1240),  # LAMAR
    (545, 1205, 715, 1270),  # FOREVER
    (240, 1380, 560, 1450),  # BRIDGE BALLROOM
    (270, 1535, 335, 1600),  # A
    (440, 1535, 505, 1600),  # B
    (1000, 1485, 1185, 1575),# ELEVATOR LOBBY
    (1185, 1440, 1300, 1495),# JUSTICE
    (1350, 1440, 1520, 1495),# GOVERNORS
    (1505, 1425, 1665, 1500),# REPRE-SENTATIVE
    (1055, 1615, 1175, 1670),# LOBBY
    (1300, 1610, 1515, 1695),# BOARDROOM FOYER
    (1185, 1625, 1280, 1868),# EXECUTIVE (vertical)
    (1545, 1660, 1610, 1840),# LIBERTY (vertical)
]
FURNITURE = [
    (275, 1025, 365, 1120),  # foyer seating clusters
    (430, 1020, 515, 1115),
    (270, 1180, 360, 1275),
    (425, 1180, 510, 1275),
    (1175, 1465, 1315, 1620),# Justice table/chairs/credenza
    (1345, 1470, 1500, 1600),# Governors
    (1505, 1455, 1670, 1610),# Representative
    (770, 540, 825, 620),    # restroom icon (upper)
    (485, 790, 540, 870),    # restroom icon (mid)
    (1025, 1725, 1150, 1810),# restroom icons (bottom)
    (125, 835, 265, 885),    # Pfluger credenza
    (300, 1648, 485, 1700),  # ballroom stage
]
BOXES = TEXT + FURNITURE

n, labels, stats, _ = cv2.connectedComponentsWithStats(mask, connectivity=8)
removed = np.zeros(n, bool)
for lbl in range(1, n):
    x, y, w, h, a = stats[lbl]
    if a < 15:  # antialiasing specks
        removed[lbl] = True
        continue
    for (x1, y1, x2, y2) in BOXES:
        if x >= x1 and y >= y1 and x + w <= x2 and y + h <= y2:
            removed[lbl] = True
            break

out = mask.copy()
out[removed[labels]] = 0

# Pixel erases: letters/icons physically touching walls (component rule can't catch them)
PIXEL_ERASE = [
    (790, 743, 851, 797),    # "RS" tail of ELEVATORS
    (1470, 1443, 1522, 1471),# "RS" tail of GOVERNORS (overlaps a wall - repaired below)
    (1042, 1738, 1074, 1843),# restroom man icon (bottom)
    (1086, 1738, 1130, 1843),# restroom woman icon (bottom)
    (1140, 1494, 1176, 1530),# stray "R" of ELEVATOR (touches wall jamb - repaired below)
    (1340, 1556, 1404, 1598),# Governors left console
    (1422, 1556, 1486, 1598),# Governors right console
    (1548, 1606, 1622, 1630),# Representative bottom console
    (1214, 1606, 1294, 1624),# Justice bottom console
    (1553, 1845, 1628, 1862),# Liberty bottom console
    (1036, 1785, 1044, 1815),# man icon arm sliver
    (1128, 1785, 1137, 1822),# woman icon arm sliver
]
for (x1, y1, x2, y2) in PIXEL_ERASE:
    out[y1:y2, x1:x2] = 0

# Repair wall segments clipped by pixel erases: (x1,y1,x2,y2,thickness)
REPAIR = [
    (1487, 1441, 1487, 1473, 7),  # vertical wall through GOVERNORS text
    (1154, 1490, 1154, 1502, 10), # wall jamb tip clipped by stray-R erase
]
for (x1, y1, x2, y2, t) in REPAIR:
    cv2.line(out, (x1, y1), (x2, y2), 1, t)

# Diagnostics: residual dark pixels inside erase boxes = letters merged into wall components
resid = []
for i, (x1, y1, x2, y2) in enumerate(BOXES):
    cnt = int(out[y1:y2, x1:x2].sum())
    if cnt > 200:
        resid.append((i, (x1, y1, x2, y2), cnt))
print("boxes with significant residual (merged-with-wall content):")
for r in resid:
    print("  ", r)

# preview: kept=black, removed=red
vis = np.full((*mask.shape, 3), 255, np.uint8)
vis[removed[labels]] = (0, 0, 255)
vis[out > 0] = (0, 0, 0)
cv2.imwrite("removed.png", vis)
cv2.imwrite("filtered.png", 255 - out * 255)
np.save("mask.npy", out)
print("kept components:", int((~removed[1:]).sum()), "/ removed:", int(removed[1:].sum()))
