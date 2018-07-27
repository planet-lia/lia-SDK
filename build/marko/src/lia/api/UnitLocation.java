package lia.api;

public class UnitLocation {
    public int id;
    public float x;
    public float y;
    public float orientation;

    public UnitLocation(int id, float x, float y, float orientation) {
        this.id = id;
        this.x = x;
        this.y = y;
        this.orientation = orientation;
    }
}
