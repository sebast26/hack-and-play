package pl.sgorecki.restclient.service;

import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import pl.sgorecki.restclient.json.AstroRecords.AstroResponse;

@Service
public class AstroService {

    private final RestClient restClient;

    public AstroService() {
        this.restClient = RestClient.create("http://api.open-notify.org");
    }

    public String getPeopleInSpace() {
        return restClient.get()
                .uri("/astros.json")
                .accept(MediaType.APPLICATION_JSON)
                .retrieve()
                .body(String.class);
    }

    public AstroResponse getAstroResponseSync() {
        return restClient.get()
                .uri("/astros.json")
                .accept(MediaType.APPLICATION_JSON)
                .retrieve()
                .body(AstroResponse.class);
    }
}
