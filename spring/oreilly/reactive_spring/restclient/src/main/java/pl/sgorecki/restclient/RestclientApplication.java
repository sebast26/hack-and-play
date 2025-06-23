package pl.sgorecki.restclient;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.support.WebClientAdapter;
import org.springframework.web.service.invoker.HttpServiceProxyFactory;
import pl.sgorecki.restclient.service.AstroInterface;

@SpringBootApplication
public class RestclientApplication {

    @Bean
    public AstroInterface astroInterface() {
        WebClient client = WebClient.create("http://api.open-notify.org/");
        HttpServiceProxyFactory factory =
                HttpServiceProxyFactory.builderFor(WebClientAdapter.create(client)).build();

        return factory.createClient(AstroInterface.class);
    }

    public static void main(String[] args) {
        SpringApplication.run(RestclientApplication.class, args);
    }

}
