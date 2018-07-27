package lia;

import lia.api.MapData;
import lia.api.StateUpdate;

public interface Callable {
    void process(MapData mapData);
    void process(StateUpdate stateUpdate, Api response);
}
