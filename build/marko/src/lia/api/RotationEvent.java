package lia.api;

public class RotationEvent {
    public EventType type;
    public int unitId;
    public Rotation rotation;

    public RotationEvent(EventType type, int unitId, Rotation rotation) {
        this.type = type;
        this.unitId = unitId;
        this.rotation = rotation;
    }
}
