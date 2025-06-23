package pl.sgorecki.restclient.json;

import java.util.List;

public class AstroRecords {
    public record Assignment(
            String name,
            String craft
    ) {
    }

    public record AstroResponse(
            String message,
            int number,
            List<Assignment> people
    ) {
    }
}
