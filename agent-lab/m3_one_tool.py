import anthropic
from tools import read_file, READ_FILE_TOOL

client = anthropic.Anthropic()

messages = [{
   "role": "user",
   "content": (
      "What does 'codex/ai/Prompt Engineering.md' say about temperature? "
      "Use the read_file tool to look it up."
   ),
}]

response = client.messages.create(
   model="claude-opus-4-8",
   max_tokens=16000,
   tools=[READ_FILE_TOOL],
   messages=messages,
)
print("stop)_reason:", response.stop_reason)

messages.append({"role": "assistant", "content": response.content})

tool_results = []
for block in response.content:
   if block.type == "tool_use":
      print(f"model wants: {block.name}({block.input})")
      output = read_file(**block.input)
      tool_results.append({
         "type": "tool_result",
         "tool_use_id": block.id,
         "content": output,
      })

messages.append({"role": "user", "content": tool_results})
final = client.messages.create(
   model="claude-opus-4-8",
   max_tokens=16000,
   tools=[READ_FILE_TOOL],
   messages=messages,
)
print(next(b.text for b in final.content if b.type == "text"))

