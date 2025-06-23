package pl.sgorecki.restclient.service;

import org.springframework.web.service.annotation.GetExchange;
import pl.sgorecki.restclient.json.AstroRecords;
import pl.sgorecki.restclient.json.AstroRecords.AstroResponse;
import reactor.core.publisher.Mono;

public interface AstroInterface {
    @GetExchange("/astros.json")
    Mono<AstroResponse> getAstroResponse();
}
