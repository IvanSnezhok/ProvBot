import telnetlib

port = '4'


cisco = telnetlib.Telnet('10.0.1.5') # Host, port, timeout
cisco.write(b'admin\r\n')
cisco.write(b'badeit\r\n')
cisco.write(b'enable\r\n')
cisco.write(b'conf t\r\n')
line = cisco.read_until(b'#', timeout=1)
print(line)
if b'am' not in line:
    cisco.write(b'logout\r\n')
    print("is not FoxGate switch")
elif '0/0/' in port:
    cisco.write(('interface ethernet ' + port + '\n').encode('ascii'))
    output1 = cisco.read_until(b'#')
    output4 = cisco.read_until(b'#')
    output2 = cisco.read_until(b'#')
else:
    cisco.write(('interface ethernet 0/0/' + port + '\n').encode('ascii'))
    output1 = cisco.read_until(b'#')
    output4 = cisco.read_until(b'#')
    output2 = cisco.read_until(b'#')
output3 = cisco.write(b'virtual-cable-test\n')
print(cisco.read_until(b'#'))
cisco.write(b'exit\n')
cisco.close()
