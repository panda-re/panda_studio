import socket

def main(message, port, IP, socketType):

    with socket.socket(socket.AF_INET, socketType) as con: 
        tries = 0
        connected = False
        while not connected and (tries < 10):
            try:
                con.connect((IP, port))
                connected = True
            except Exception as e:
                tries = tries + 1
                pass #Do nothing, just try again

        if connected:
            con.send(message)
            #response = con.recv(32768)
            con.close()
        else:
            raise("Connection_Error")

if __name__ == "__main__":
    main()