package lia.api;

public class Unit {
    public int id;
    public int health;
    public float x;
    public float y;
    public float orientation;
    public ThrustSpeed thrustSpeed;
    public Rotation rotation;
    public boolean canShoot;
    public int nBullets;
    public OpponentInView[] opponentsInView;
    public BulletInView[] bulletsInView;

    public Unit(int id,
                int health,
                float x, float y,
                float orientation,
                ThrustSpeed thrustSpeed,
                Rotation rotation,
                boolean canShoot,
                int nBullets,
                OpponentInView[] opponentsInView,
                BulletInView[] bulletsInView) {
        this.id = id;
        this.health = health;
        this.x = x;
        this.y = y;
        this.orientation = orientation;
        this.thrustSpeed = thrustSpeed;
        this.rotation = rotation;
        this.canShoot = canShoot;
        this.nBullets = nBullets;
        this.opponentsInView = opponentsInView;
        this.bulletsInView = bulletsInView;
    }
}

