package lia.api;

public class ShootEvent {
    public EventType type;
    public int unitId;

    public ShootEvent(EventType type, int unitId) {
        this.type = type;
        this.unitId = unitId;
    }
}
