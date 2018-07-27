package logic;

import com.adamldavis.pathfinder.AntPathFinder;
import com.adamldavis.pathfinder.PathGrid;
import lia.Api;
import lia.api.Unit;
import lia.api.Rotation;
import lia.api.ThrustSpeed;

import java.util.ArrayList;

public class PathFollower {

    private static final float ALLOWED_UNIT_OFFSET = 2f;
    private static final float ALLOWED_ANGLE_OFFSET = 15f;

    private ArrayList<Vector2> points;
    private int nextPointIndex = 0;
    private Vector2 vector2 = new Vector2();

    /**
     * Find and create a path from (x1, y1) to (x2, y2) on specified grid.
     * */
    public PathFollower(PathGrid grid, int x1, int y1, int x2, int y2) {

        AntPathFinder finder = new AntPathFinder(80);

        int[] moves = finder.findPath(grid, x1, y1, x2, y2);
        if (moves == null) return;

        // Convert moves to points
        points = new ArrayList<>(moves.length + 1);

        int x = x1;
        int y = y1;

        for (int move : moves) {
            switch (move) {
                case 0 : y -= 1; break;
                case 1 : x += 1; break;
                case 2 : y += 1; break;
                case 3 : x -= 1; break;
            }
            points.add(new Vector2(x, y));
        }

        optimizePathPoints();
    }

    private void optimizePathPoints() {
        // Removes redundant points on vertical and horizontal lines
        for (int i = 1; i < points.size() - 1; i++) {
            Vector2 p1 = points.get(i - 1);
            Vector2 p2 = points.get(i);
            Vector2 p3 = points.get(i + 1);

            if (p1.x == p2.x && p2.x == p3.x ||
                    p1.y == p2.y && p2.y == p3.y) {
                points.remove(i);
                i--;
            }
        }
        // TODO Improvements:
        //  1: Remove redundant points on diagonal lines
        //  2: Find shortcuts by choosing points further
        //     apart and checking if there are no obstacles
        //     on the line between them.
        //  3: Make units move in circular motions without
        //     stopping.
    }

    /**
     * Follows the previously chosen path and writes
     * the needed step to api.
     * @return true if the unit is at the last point.
     * */
    public boolean follow(Unit unit, Api api) {
        int x = (int) unit.x;
        int y = (int) unit.y;

        Vector2 nextPoint;

        // If unit is close enough to the nextPoint then
        // move to the next point
        while (true) {
            if (nextPointIndex >= points.size()) return true;

            int oldIndex = nextPointIndex;
            nextPoint = points.get(nextPointIndex);

            if (Math.abs(nextPoint.x - x) < ALLOWED_UNIT_OFFSET &&
                    Math.abs(nextPoint.y - y) < ALLOWED_UNIT_OFFSET) {

                nextPointIndex++;

                if (nextPointIndex == points.size()) {
                    // No more points to visit
                    return true;
                }
            }
            if (oldIndex == nextPointIndex) {
                break;
            }
        }

        // Get angle between current point and nextPoint
        vector2.set(nextPoint);
        vector2.sub(x, y);
        float angle = unit.orientation - vector2.angle();
        if (angle > 180) angle -= 360;
        else if (angle < -180) angle += 360;

        // If the angle is small enough are close move forward
        if (Math.abs(angle) < ALLOWED_ANGLE_OFFSET) {
            api.setRotationSpeed(unit.id, Rotation.NONE);
            api.setThrustSpeed(unit.id, ThrustSpeed.FORWARD);
        }
        // Else rotate to the needed direction
        else {
            api.setThrustSpeed(unit.id, ThrustSpeed.NONE);
            if (angle < 0f) {
                api.setRotationSpeed(unit.id, Rotation.LEFT);
            } else {
                api.setRotationSpeed(unit.id, Rotation.RIGHT);
            }
        }

        return false;
    }

    public void printPoints() {
        for (Vector2 point : points) System.out.print(point);
        System.out.println();
    }
}
