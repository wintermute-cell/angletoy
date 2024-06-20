import meshlib.mrmeshpy as mmp

path = ".../map_thresh_whitewalls.png"
dm = mmp.loadDistanceMapFromImage(mmp.Path(path), 0)


pl2 = mmp.distanceMapTo2DIsoPolyline(dm, isoValue=127)
mesh = mmp.triangulateContours(pl2.contours2())

mmp.saveMesh(mesh, mmp.Path(".../mesh.obj"))

if __name__ == "__main__":
    pass
