word_size = 8 # 64bit only for now

def get_reg(name):
  for req in registers().Regs:
    if req.Name == name:
      return int(req.Value, 16)

def read_word(addr):
  v = 0
  m = examine_memory(addr, word_size)
  s = 1
  for b in m.Mem:
    v += b * s
    s = s * 256
  return v

def hex(d, n = 0):
  if d == 0:
    return lpad("0", n, "0")

  lookup = {10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "e", 16: "f"}
  s = ""
  while d > 0:
    r = d % 16
    d = d // 16
    if r >= 10:
      r = lookup[r]
    else:
      r = str(r)
    s = r + s
  return lpad(s, n, "0")

def lpad(s, n, c = " "):
  while len(s) < n:
    s = c + s
  return s

def rpad(s, n, c = " "):
  while len(s) < n:
    s = s + c
  return s

def ascii_table(rows, align_right = {}):
  widths = []
  for row in rows:
    for i, col in enumerate(row):
      if len(widths) < i+1:
        widths.append(0)
      widths[i] = max(widths[i], len(col))

  s = ""
  for row in rows:
    for i, col in enumerate(row):
      width = widths[i]
      if align_right.get(i, False):
        col = lpad(col, width)
      else:
        col = rpad(col, width)
      s += col + "  "
    s += "\n"
  return s

def getg():
  # TODO(fg) there is probably a better way to implement this.
  g = raw_command("goroutine").State.SelectedGoroutine
  for gp in eval(None, "runtime.allgs").Variable.Value:
    if gp.goid == g.ID:
      return gp

def stack():
  g = getg()
  bp = get_reg("Rbp")
  ip = get_reg("Rip")
  sp = get_reg("Rsp")
  regs = {"sp": sp, "bp": bp, "ip": ip}

  addr_list = []
  addr_dict = {}
  offset = 0
  while True:
    addr = g.stack.hi+offset-word_size
    addr_info = {
      "addr": addr,
      "val": read_word(addr),
      "offset": offset,
      "regs": [],
      "note": [],
      "func": None,
      "arg": None,
      "local": None,
      "fp": False,
    }
    for (name, val) in regs.items():
      if addr == val:
        addr_info["regs"].append(name)
        
    addr_list.append(addr_info)
    addr_dict[addr] = addr_info
    offset -= word_size
    if addr <= sp:
      break
  
  for f in stacktrace(g.goid, 128, True).Locations:
    fp_addr = g.stack.hi+f.FramePointerOffset
    if fp_addr > 0 and len(addr_dict[fp_addr]["note"]) == 0:
      addr_dict[fp_addr]["note"].append("frame pointer for "+f.Function.Name_)
      pc_addr = fp_addr+word_size
      pc = read_word(pc_addr)
      ins = disassemble(None, pc, pc+1).Disassemble[0]
      addr_dict[pc_addr]["note"].append("return addr to "+ins.Loc.Function.Name_)

    for arg in f.Arguments:
      addr_dict[arg.Addr]["note"].append("arg "+arg.Name+" "+arg.Type)
    for local in f.Locals:
      addr = local.Addr // 8 * 8
      if addr_dict.get(addr):
        addr_dict[addr]["note"].append("var "+local.Name+" "+local.Type)

  return addr_list

# stackannotate (alias sa) will print an annotated stack dump.
def command_stackannotate():
  rows = [["regs", "addr", "offset", "value", "explanation"]]
  for addr_info in stack():
    regs = ""
    if len(addr_info["regs"]) > 0:
      regs = ",".join(addr_info["regs"])+" -->"

    note = "?"
    if len(addr_info["note"]) > 0:
      note = ", ".join(addr_info["note"])
      
    rows.append([
      regs,
      hex(addr_info["addr"]),
      lpad(str(addr_info["offset"]), 6),
      lpad(hex(addr_info["val"]), word_size*2+2),
      note,
    ])
  print(ascii_table(rows))

def main():
  dlv_command("config alias stackannotate sa")
