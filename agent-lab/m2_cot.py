import json
import anthropic

client = anthropic.Anthropic()

QUESTION = "If a debt of $2000 grows at 24% APR and I pay $100/month, roughly how long until it is gone?  Give the number of months."

SCHEMA = {
   "type": "object",
   "properties": {
      "reasoning": {"type": "string", "description": "Step-by-step working"},
      "months": {"type": "integer", "description": "Whole months to payoff"},
   },
   "required": ["reasoning", "months"],
   "additionalProperties": False,
}

response = client.messages.create(
   model="claude-opus-4-8",
   max_tokens=16000,
   output_config={"format": {"type": "json_schema", "schema": SCHEMA}},
   messages=[{
      "role": "user",
      "content": (
         "Think through this step by step in the 'reasoning' field, "
         "then put the final whole-month count in 'months',\n\n"
         f"{QUESTION}"
      ),
   }],
)

data = json.loads(next(b.text for b in response.content if b.type == "text"))
print("reasoning:", data["reasoning"])
print("months:", data["months"])


