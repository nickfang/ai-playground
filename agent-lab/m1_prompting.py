import anthropic

client = anthropic.Anthropic()

PAGE = """Avalanche vs Snowball are two debt-payoff strategies.  Avalanche pays the highest-interest debt first (mathematically optimal). Snowball pays the smalles balance first (psychologically motivating)."""

PAGE2 = """A digest of a YouTube Channel: Hi, welcome to my youtube channel.  Today we are going to talk about rabbits...""" 

def classify_zero_shot(page_body: str) -> str:
   prompt = (
      "Classify this knowledge-base page as exactly one of: "
      "source, entity, concept, topic, overview, project. "
      "Reply with only the label.\n\n"
      f"Page:\n{page_body}"
   )
   response = client.messages.create(
      model="claude-opus-4-8",
      max_tokens=16,
      messages=[{"role": "user", "content": prompt}],
   )
   return response.content[0].text.strip()

print("zero-shot:", classify_zero_shot(PAGE))

def classify_few_shot(page_body: str) -> str:
   prompt = (
      "Classify each knowledge-base page as exactly one of: "
      "source, entity, concept, topic, overview, project.\n\n"
      "Page: A digest of a YouTube video on CI/CD.\nLabel: source\n\n"
      "Page: Chase Sapphire Preferred, a travel credit card.\nLabel: entity\n\n"
      "page: Compound interest, the idea that interest earns interest.\nLabel: concept\n\n"
      f"Page: {page_body}\nLabel:"
   )
   response = client.messages.create(
      model="claude-opus-4-8",
      max_tokens=16,
      messages=[{"role": "user", "content": prompt}],
   )
   return response.content[0].text.strip()

print("few-shot: ", classify_few_shot(PAGE))

