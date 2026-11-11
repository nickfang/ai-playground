from langchain.agents import create_agent
from lc_tools import LC_TOOLS

SYSTEM = (
   "You are a research assistant for a personal knowledge vault (Cortex). "
   "Answer the user's question using ONLY what you find in the vault"
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to.  "
   "Cite the file paths you drew from.  If the vault does not cover it, say so."
)

from langchain_anthropic import ChatAnthropic
from langchain_core.messages import SystemMessage, HumanMessage, ToolMessage
from lc_tools import LC_TOOLS, TOOL_BY_NAME

SYSTEM = (
   "You are a research assistant for a personal knowledge vault (Cortex). "
   "Answer the user's question using ONLY what you find in the vault"
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to.  "
   "Cite the file paths you drew from.  If the vault does not cover it, say so."
)

model = ChatAnthropic(model="claude-opus-4-8", max_tokens=16000)
model_with_tools = model.bind_tools(LC_TOOLS)

def run_agent(question: str, max_turns: int = 12) -> str:
   messages = [SystemMessage(SYSTEM), HumanMessage(question)]
   for turn in range(max_turns):
      ai = model_with_tools.invoke(messages)
      messages.append(ai)
      if not ai.tool_calls:
         return ai.text
      for call in ai.tool_calls:
         print(f"   [turn {turn}] {call['name']}({call['args']})")
         output = TOOL_BY_NAME[call["name"]].invoke(call["args"])
         messages.append(ToolMessage(content=output, tool_call_id=call["id"]))
   return "Stopped: hit the max-turns limit without finishing."

if __name__ == "__main__":
   print(run_agent("What does the vault say about garlic?"))
