package pl.sgorecki.koog

import ai.koog.agents.core.agent.AIAgent
import ai.koog.agents.core.tools.ToolRegistry
import ai.koog.agents.core.tools.annotations.LLMDescription
import ai.koog.agents.core.tools.annotations.Tool
import ai.koog.agents.core.tools.reflect.ToolSet
import ai.koog.agents.core.tools.reflect.tools
import ai.koog.agents.features.eventHandler.feature.handleEvents
import ai.koog.prompt.executor.clients.openai.OpenAIModels
import ai.koog.prompt.executor.llms.all.simpleOpenAIExecutor
import ai.koog.prompt.llm.LLModel
import kotlinx.serialization.Serializable

@Serializable
data class Contact(val id: Int, val name: String, val lastName: String, val phone: String)

val contactList = listOf(
    Contact(100, "Alice", "Smith", "+1 415 111 1111"),
    Contact(101, "Bob", "Johnosn", "+38 1123 123 123"),
    Contact(103, "Daniel", "Anderson", "+50 123 123 123"),
    Contact(104, "Daniel", "Garcia", "+11 991 123 123")
)

val contactMap = contactList.associateBy { it.id }

class MoneyTransferTools : ToolSet {


    @Tool
    @LLMDescription("Sends money to specified recipient with given amount and purpose")
    fun sendMoney(
        @LLMDescription("The ID of the user initiating the transfer")
        senderId: Int,
        @LLMDescription("Amount of money to send in euros")
        amount: Double,
        @LLMDescription("Name of the recipient")
        recipientID: Int,
        @LLMDescription("Purpose of the money transfer")
        purpose: String
    ): String {
        val recipient = contactMap[recipientID] ?: return "Invalid recipient"
        println("-----------")
        println("Sending money to ${recipient.name} ${recipient.lastName} ${recipient.phone} for $amount euros with purpose: $purpose")
        println("Please confirm the transaction by typing 'yes'")
        println("-----------")
        val confirmation = readln()
        return if (confirmation.lowercase() == "yes") "Money sent" else "Transaction declined"
    }

    @Tool
    @LLMDescription("Returns a list of all contacts")
    fun getContacts(
        @LLMDescription("The unique identifier of the user whose contact list is being retrieved.")
        userId: Int
    ): String {
        return contactList.joinToString("\n") { "${it.id} - ${it.name} ${it.lastName}" }
    }

    @Tool
    @LLMDescription("Helps to identify the correct recipient when multiple contacts have similar names")
    fun chooseRecipient(
        @LLMDescription("The unique identifier of the user who initiated the transfer")
        userId: Int,
        @LLMDescription("The ambiguous name that needs to be clarified")
        confusingName: String
    ): String {
        val matches = contactList.filter {
            it.name.contains(confusingName, ignoreCase = true) ||
                    it.lastName.contains(confusingName, ignoreCase = true)
        }

        if (matches.isEmpty()) return "No matching contacts found"
        if (matches.size == 1) return matches[0].id.toString()

        println("Multiple matches found. Please select the correct recipient:")
        matches.forEachIndexed { index, contact ->
            println("$index) ${contact.name} ${contact.lastName} (ID: ${contact.id})")
        }

        println("Enter the number of the correct recipient: ")
        val selection = readln().toIntOrNull()

        return when(selection) {
            null -> "Invalid selection"
            !in matches.indices -> "Invalid selection"
            else -> matches[selection].id.toString()
        }
    }
}

suspend fun main() {
    val executor = simpleOpenAIExecutor(System.getenv("OPENAI_API_KEY"))
    val model: LLModel = OpenAIModels.CostOptimized.GPT4_1Mini

    val toolRegistry = ToolRegistry {
        tools(MoneyTransferTools())
    }
    val agent = AIAgent(
        executor = executor,
        llmModel = model,
        toolRegistry = toolRegistry,
        systemPrompt = "You're a banking assistant. Accompany the user with their request."
    ) {
        handleEvents {
            onBeforeLLMCall { ctx ->
                println("Request to LLM:")
                println("    # Messages:")
                ctx.prompt.messages.forEach { println("      - $it") }
                println("    # Tools:")
                ctx.tools.forEach { println("      - $it") }
            }
            onAfterLLMCall { ctx ->
                println("Response from LLM:")
                ctx.responses.forEach { println("      - $it") }
            }
        }
    }
    val userMessage = "Send 25 euros to Daniel for dinner at restaurant."
    val result = agent.run(userMessage)
    println("Final result: $result")
}