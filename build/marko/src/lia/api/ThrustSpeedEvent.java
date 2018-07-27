package lia.api;

public class ThrustSpeedEvent {
    public EventType type;
    public int unitId;
    public ThrustSpeed speed;

    public ThrustSpeedEvent(EventType type, int unitId, ThrustSpeed speed) {
        this.type = type;
        this.unitId = unitId;
        this.speed = speed;
    }
}
