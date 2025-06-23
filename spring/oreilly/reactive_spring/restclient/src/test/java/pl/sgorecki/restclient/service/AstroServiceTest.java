package pl.sgorecki.restclient.service;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import pl.sgorecki.restclient.json.AstroRecords.AstroResponse;

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
}