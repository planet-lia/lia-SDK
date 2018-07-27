package lia;

import com.google.gson.Gson;
import lia.api.*;

import java.util.ArrayList;

/**
 * Used for building a response message that is later
 * sent to the game engine.
 **/
public class Api {

    private long uid;

    private ArrayList<ThrustSpeedEvent> thrustSpeedEvents;
    private ArrayList<RotationEvent> rotationEvents;
    private ArrayList<ShootEvent> shootEvents;

    protected Api() {
        thrustSpeedEvents = new ArrayList<>();
        rotationEvents = new ArrayList<>();
        shootEvents = new ArrayList<>();
    }

    protected void setUid(long uid) {
        this.uid = uid;
    }

    /** Change thrust speed of a unit */
    public void setThrustSpeed(int unitId, ThrustSpeed speed) {
        thrustSpeedEvents.add(
                new ThrustSpeedEvent(EventType.SET_THRUST_SPEED, unitId, speed)
        );
    }

    /** Change rotation speed of a unit */
    public void setRotationSpeed(int unitId, Rotation rotation) {
        rotationEvents.add(
                new RotationEvent(EventType.SET_ROTATION, unitId, rotation)
        );
    }

    /** Make a unit shoot */
    public void shoot(int unitId) {
        shootEvents.add(
                new ShootEvent(EventType.SHOOT, unitId)
        );
    }

    protected String toJson() {
        ThrustSpeedEvent[] thrust = new ThrustSpeedEvent[thrustSpeedEvents.size()];
        thrust = thrustSpeedEvents.toArray(thrust);

        RotationEvent[] rotation = new RotationEvent[rotationEvents.size()];
        rotation = rotationEvents.toArray(rotation);

        ShootEvent[] shoot = new ShootEvent[shootEvents.size()];
        shoot = shootEvents.toArray(shoot);

        Response response =  new Response(
                uid,
                MessageType.RESPONSE,
                thrust,
                rotation,
                shoot
        );
        return (new Gson()).toJson(response);
    }
}