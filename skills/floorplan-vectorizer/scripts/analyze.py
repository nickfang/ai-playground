import cv2
import numpy as np

img = cv2.imread("source.png")
gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
# dark linework mask (walls, text, furniture are near-black; fills are orange ~mid; bg cream ~bright)
mask = (gray < 110).astype(np.uint8) * 255

n, labels, stats, centroids = cv2.connectedComponentsWithStats(mask, connectivity=8)
areas = stats[1:, cv2.CC_STAT_AREA]
order = np.argsort(areas)[::-1]
print(f"total components: {n-1}")
print("top 15 by area:")
for i in order[:15]:
    lbl = i + 1
    x, y, w, h, a = stats[lbl]
    print(f"  label={lbl} area={a} bbox=({x},{y},{w}x{h})")

# visualize: largest component in black, everything else in red
out = np.full((*gray.shape, 3), 255, np.uint8)
biggest = order[0] + 1
out[(labels > 0) & (labels != biggest)] = (0, 0, 255)
out[labels == biggest] = (0, 0, 0)
cv2.imwrite("components.png", out)
print("wrote components.png (largest=black, rest=red)")
