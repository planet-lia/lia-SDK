package logic;

import static logic.Constants.DEGREES_TO_RADIANS;
import static logic.Constants.RADIANS_TO_DEGREES;

public class Vector2 {

    public float x;
    public float y;

    public Vector2() {
        this.x = 0;
        this.y = 0;
    }

    public Vector2(float x, float y) {
        this.x = x;
        this.y = y;
    }


    public void set(float x, float y) {
        this.x = x;
        this.y = y;
    }

    public void set(Vector2 v) {
        set(v.x, v.y);
    }

    /** @return the angle in degrees of this vector (point) relative to the x-axis. Angles are towards the positive y-axis (typically
     * counter-clockwise) and between 0 and 360.
     */
    public float angle() {
        double angle = Math.atan2(y, x) * RADIANS_TO_DEGREES;
        if (angle < 0) angle += 360f;
        return (float) angle;
    }

    public Vector2 add(Vector2 v) {
        x += v.x;
        y += v.y;
        return this;
    }

    public Vector2 sub(Vector2 v) {
        return v;
    }

    public Vector2 sub(float x, float y) {
        this.x -= x;
        this.y -= y;
        return this;
    }

    /** Rotates the Vector2 by the given angle, counter-clockwise assuming the y-axis points up.
     * @param degrees the angle in degrees
     */
    public Vector2 rotate(float degrees) {
        return rotateRad(degrees * DEGREES_TO_RADIANS);
    }


    /** Rotates the Vector2 by the given angle, counter-clockwise assuming the y-axis points up.
     * @param radians the angle in radians
     */
    public Vector2 rotateRad(Float radians) {
        double cos = Math.cos(radians);
        double sin = Math.sin(radians);

        double newX = this.x * cos - this.y * sin;
        double newY = this.x * sin + this.y * cos;

        this.x = (float) newX;
        this.y = (float) newY;

        return this;
    }

    @Override
    public String toString() {
        return "(" + x + "," + y + ")";
    }

}
