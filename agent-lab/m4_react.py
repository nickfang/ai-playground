import anthropic
from tools import TOOLS, dispatch

client = anthropic.Anthropic()

SYSTEM = (
   "You are an instructor that teaches information from the personal knowledge vault (Cortex). "
   "Answer the user's question using ONLY what you find in the vault. "
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to. "
   "Cite the file paths you drew from.  If the vault does not cover it, say so. "
)

def run_agent(question: str, max_turns: int = 12) -> str:
   messages = [{"role": "user", "content": question}]
   for turn in range(max_turns):
      response = client.messages.create(
         model="claude-opus-4-8",
         max_tokens=16000,
         system=SYSTEM,
         tools=TOOLS,
         messages=messages,
      )
      print(f"   [turn {turn}] tokens: in={response.usage.input_tokens} out={response.usage.output_tokens}")
      
      if response.stop_reason == "end_turn":
         return next(b.text for b in response.content if b.type == "text")

      messages.append({"role": "assistant", "content": response.content})

      tool_results = []
      for block in response.content:
         if block.type == "tool_use":
            print(f"  [turn {turn}] {block.name}({block.input})")
            output = dispatch(block.name, block.input)
            tool_results.append({
               "type": "tool_result",
               "tool_use_id": block.id,
               "content": output,
            })
      messages.append({"role": "user", "content": tool_results})

   return "Stopped: hit the max-turns limit without finishing."

if __name__ == "__main__":
   answer = run_agent("What does the vault say about fishing?")
   print("\n=== Answer ===")
   print(answer)
