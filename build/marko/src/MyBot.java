import com.adamldavis.pathfinder.PathGrid;
import lia.Api;
import lia.Callable;
import lia.NetworkingClient;
import lia.api.*;
import logic.UnitData;
import logic.PathFinding;
import logic.Vector2;
import java.util.HashMap;


/**
 * Place to write the logic for your bots.
 * */
public class MyBot implements Callable {

    private HashMap<Integer, UnitData> unitsData = new HashMap<>();

    private Vector2 bottomRightCorner = new Vector2(140, 2);
    private Vector2 topLeftCorner = new Vector2(2, 78);

    /** Called only once when the game is initialized. */
    @Override
    public synchronized void process(MapData mapData) {
            // Convert map with obstacles to a grid that will be used for
            // path finding algorithm
            PathGrid grid = PathFinding.createGrid(
                    (int) mapData.width,
                    (int) mapData.height,
                    mapData.obstacles
            );

            // Store data related to each unit in UnitData object.
            // All UnitData objects will be accessible through unitsData map.
            for (UnitLocation unit : mapData.unitLocations) {
                UnitData data = new UnitData(unit.id, grid);
                unitsData.put(unit.id, data);
            }
    }

    /** Repeatedly called from game engine with game state updates.  */
    @Override
    public synchronized void process(StateUpdate stateUpdate, Api api) {
        // Go through all units
        for (Unit unit : stateUpdate.units) {
            // Get data for this unit
            UnitData unitData = unitsData.get(unit.id);

            // If the unit is not following any path, set it
            if (unitData.pathNotSet) {
                if (unitData.goingToTopLeftCorner) {
                    unitData.setPathToFollow(unit, topLeftCorner.x,  topLeftCorner.y);
                } else {
                    unitData.setPathToFollow(unit, bottomRightCorner.x,  bottomRightCorner.y);
                }
                unitData.goingToTopLeftCorner = !unitData.goingToTopLeftCorner;
            }

            // Make the unit follow the path
            unitData.followPath(unit, api);

            // Shoot if you see an opponent and your gun is loaded
            if (unit.opponentsInView.length > 0 && unit.canShoot)
                api.shoot(unit.id);
        }
    }

    public static void main(String[] args) throws Exception {
        NetworkingClient.connectNew(args, new MyBot());
    }
}
