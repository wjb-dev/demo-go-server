#include <gtest/gtest.h>
#include "generated/service.pb.h"
#include "generated/service.grpc.pb.h"

using v1::EchoRequest;
using v1::EchoResponse;
using v1::EchoService;

class StubService : public EchoService::Service {
 public:
  grpc::Status Echo(grpc::ServerContext*, const EchoRequest* req,
                    EchoResponse* resp) override {
    resp->set_message(req->message());
    return grpc::Status::OK;
  }
};

TEST(EchoTest, EchoReturnsSameMessage) {
  StubService stub;
  EchoRequest req;
  req.set_message("hello");
  EchoResponse resp;
  grpc::ServerContext ctx;

  grpc::Status status = stub.Echo(&ctx, &req, &resp);
  EXPECT_TRUE(status.ok());
  EXPECT_EQ(resp.message(), "hello");
}
