package logic;

import com.adamldavis.pathfinder.PathGrid;
import com.adamldavis.pathfinder.SimplePathGrid;
import lia.api.Obstacle;

public class PathFinding {

    private static final int OFFSET_FROM_OBSTACLE = 1;


    public static PathGrid createGrid(int mapWidth, int mapHeight, Obstacle[] obstacles) {
        PathGrid grid = new SimplePathGrid(mapWidth, mapHeight);

        // Set grid to false where we don't want for units to move (based on the
        // positions of obstacles)
        for (Obstacle obstacle : obstacles) {
            int x1 = (int) obstacle.x - OFFSET_FROM_OBSTACLE;
            int x2 = (int) (obstacle.x + obstacle.width) + OFFSET_FROM_OBSTACLE;

            for (int x = x1; x < x2; x++) {
                int y1 = (int) obstacle.y - OFFSET_FROM_OBSTACLE;
                int y2 = (int) (obstacle.y + obstacle.height) + OFFSET_FROM_OBSTACLE;

                for (int y = y1; y < y2; y++) {
                    // If x and y are in obstacle or very close to it, set grid to true
                    if (x < grid.getWidth() && x >= 0 && y < grid.getHeight() && y >= 0) {
                        grid.setGrid(x, y, true);
                    }
                }
            }
        }

        return grid;
    }
}
