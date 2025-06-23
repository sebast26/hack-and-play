package pl.sgorecki.restclient.service;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import pl.sgorecki.restclient.json.AstroRecords.AstroResponse;
import reactor.test.StepVerifier;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
class AstroServiceTest {
    @Autowired
    private AstroService astroService;

    @Test
    void getPeopleInSpace() {
        String people = astroService.getPeopleInSpace();
        assertNotNull(people);
        assertTrue(people.contains("success"));
        System.out.println(people);
    }

    @Test
    void getAstroResponseSync() {
        AstroResponse response = astroService.getAstroResponseSync();
        assertNotNull(response);
        assertEquals("success", response.message());
        assertTrue(response.number() >= 0);
        assertEquals(response.number(), response.people().size());
        System.out.println(response);
    }

    @Test
    void getAstroResponseAsync() {
        AstroResponse response = astroService.getAstroResponseAsync()
                .block(Duration.ofSeconds(10));
        assertNotNull(response);
        assertEquals("success", response.message());
        assertTrue(response.number() >= 0);
        assertEquals(response.number(), response.people().size());
        System.out.println(response);
    }

    @Test
    void getAstroResponseAsyncStepVerifier() {
        astroService.getAstroResponseAsync()
                .as(StepVerifier::create)
                .assertNext(response -> {
                    assertNotNull(response);
                    assertEquals("success", response.message());
                    assertTrue(response.number() >= 0);
                    assertEquals(response.number(), response.people().size());
                    System.out.println(response);
                })
                .verifyComplete();
    }

    @Test
    void getAstroResponseFromInterface(@Autowired AstroInterface astroInterface) {
        AstroResponse response = astroInterface.getAstroResponse()
                .block(Duration.ofSeconds(10));
        assertNotNull(response);
        assertAll(
                () -> assertEquals("success", response.message()),
                () -> assertTrue(response.number() >= 0),
                () -> assertEquals(response.number(), response.people().size())
        );
        System.out.println(response);
    }
}