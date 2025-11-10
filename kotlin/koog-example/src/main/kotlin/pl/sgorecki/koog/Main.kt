package pl.sgorecki.koog

import ai.koog.prompt.dsl.prompt
import ai.koog.prompt.executor.clients.LLMClient
import ai.koog.prompt.executor.clients.openai.OpenAILLMClient
import ai.koog.prompt.executor.clients.openai.OpenAIModels
import ai.koog.prompt.llm.LLModel
import ai.koog.prompt.params.LLMParams

suspend fun main() {
    val client: LLMClient = OpenAILLMClient(System.getenv("OPENAI_API_KEY"))
    val model: LLModel = OpenAIModels.CostOptimized.GPT4_1Mini

    val userMessage = readln()
    val prompt = prompt(
        id = "translation-request",
        params = LLMParams(temperature = 0.7)
    ) {
        user("Translate this sentence to German: $userMessage")
    }
    val responses = client.execute(prompt, model)
    with(responses.first()) {
        println("Translation: $content")
        println("Meta info: $metaInfo")
    }
}