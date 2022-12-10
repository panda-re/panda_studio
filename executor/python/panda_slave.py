import grpc
import concurrent.futures as futures

import panda_interface_pb2_grpc as pb_grpc
import panda_interface_pb2 as pb

PORT = '[::]:50051'

class PandaExecutorServicer(pb_grpc.PandaExecutorServicer):
    def BootMachine(self, request, context):
        return super().BootMachine(request, context)
    
    def RunCommand(self, request, context):
        print(request.command)
        return pb.RunCommandReply(statusCode=32)
        # return super().RunCommand(request, context)

def serve():
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
  pb_grpc.add_PandaExecutorServicer_to_server(
      PandaExecutorServicer(), server)
  server.add_insecure_port(PORT)
  print(f'panda agent grpc server listening on port {PORT}')
  server.start()
  server.wait_for_termination()

if __name__ == "__main__":
    serve()