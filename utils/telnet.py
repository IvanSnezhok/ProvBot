import telnetlib

HOST = '10.0.1.5'
user = 'admin'
password = 'badeit'


def to_bytes(line):
    return bytes(line, "UTF-8") + b'/r'


def telnet_cable_test(host: str, port: str):
    with telnetlib.Telnet(host) as tn:
        tn.write(b"admin\r")
        tn.write(b'badeit\r')
        tn.read_until(b'#')
        result = {}
        if tn.read_until(b'#'):
            tn.write(to_bytes("enable"))
            tn.write(b'conf t\r\n')
            tn.read_until(b'#')
            command = 'interface ethernet 1/{}'.format(port)
            command_1 = 'show am interface ethernet 1/{}'.format(port)
            command_2 = 'show mac-address-table interface ethernet 1/{}'.format(port)
            command_3 = ' '.format(port)
            command_4 = 'virtual-cable-test'
            commands = [command, command_1, command_2, command_3, command_4]
            for command in commands:
                tn.write(to_bytes(command))
                output = tn.read_until(b"#", timeout=5).decode('utf-8')
                result[command] = output.replace('\r\n', '\n')
            tn.write(to_bytes('exit'))
            tn.write(to_bytes('exit'))
            tn.write(to_bytes('exit'))
            return result
        elif tn.read_until(b'>'):
            tn.write(to_bytes("enable"))
            tn.write(b'conf t\r\n')
            tn.read_until(b'>')
            command = 'interface ethernet 0/0/{}'.format(port)
            command_1 = 'show am interface ethernet 0/0/{}'.format(port)
            command_2 = 'show mac-address-table interface ethernet 0/0/{}'.format(port)
            command_3 = 'show interface ethernet 0/0/{} detail'.format(port)
            command_4 = 'virtual-cable-test'
            commands = [command, command_1, command_2, command_3, command_4]
            for command in commands:
                tn.write(to_bytes(command))
                output = tn.read_until(b">", timeout=5).decode('utf-8')
                result[command] = output.replace('\r\n', '\n')
            tn.write(to_bytes('exit'))
            tn.write(to_bytes('exit'))
            tn.write(to_bytes('exit'))
            return result
        else:
            print("Error")
            return "Error"


if __name__ == '__main__':
    print(telnet_cable_test(HOST, '4'))
