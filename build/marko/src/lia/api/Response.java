package lia.api;

public class Response {
    public long uid;
    public MessageType type;
    public ThrustSpeedEvent[] thrustSpeedEvents;
    public RotationEvent[] rotationEvents;
    public ShootEvent[] shootEvents;

    public Response(long uid, MessageType type,
                    ThrustSpeedEvent[] thrustSpeedEvents,
                    RotationEvent[] rotationEvents,
                    ShootEvent[] shootEvents) {
        this.uid = uid;
        this.type = type;
        this.thrustSpeedEvents = thrustSpeedEvents;
        this.rotationEvents = rotationEvents;
        this.shootEvents = shootEvents;
    }
}
