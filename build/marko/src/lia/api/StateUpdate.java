package lia.api;

public class StateUpdate {
    public long uid;
    public MessageType type;
    public float time;
    public Unit[] units;

    public StateUpdate(long uid, MessageType type, float time, Unit[] units) {
        this.uid = uid;
        this.type = type;
        this.time = time;
        this.units = units;
    }
}


