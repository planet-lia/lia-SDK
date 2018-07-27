package logic;

import com.adamldavis.pathfinder.PathGrid;
import lia.Api;
import lia.api.Unit;

public class UnitData {

    public int id;
    public PathFollower follower;
    public PathGrid grid;
    public Vector2 position;
    public boolean pathNotSet;
    public boolean goingToTopLeftCorner;

    public UnitData(int id, PathGrid grid) {
        this.id = id;
        this.grid = grid;
        this.position = new Vector2();
        this.pathNotSet = true;
    }

    public void setPathToFollow(Unit unit, float x, float y) {
        follower = new PathFollower(grid, (int) unit.x, (int) unit.y, (int) x, (int) y);
        pathNotSet = false;
    }

    public void followPath(Unit unit, Api api) {
        if (follower == null) return;
        pathNotSet = follower.follow(unit, api);
    }
}
