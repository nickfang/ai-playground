import anthropic

client = anthropic.Anthropic()

response = client.messages.create(
   model="claude-opus-4-8",
   max_tokens=1024,
   messages=[{"role": "user", "content": "In one sentence, what is a ReAct agent?"}],
)

for block in response.content:
   if block.type == "text":
      print(block.text)
