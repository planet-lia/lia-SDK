package lia.api;

public class BulletInView {
    float x;
    float y;
    float orientation;
    float velocity;

    public BulletInView(float x, float y,
                        float orientation, float velocity) {
        this.x = x;
        this.y = y;
        this.orientation = orientation;
        this.velocity = velocity;
    }
}
