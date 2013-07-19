package protocol

import "testing"

var (
	emptyMetadataResponse = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}

	brokersNoTopicsMetadataResponse = []byte{
		0x00, 0x00, 0x00, 0x02,

		0x00, 0x00, 0xab, 0xff,
		0x00, 0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't',
		0x00, 0x00, 0x00, 0x33,

		0x00, 0x01, 0x02, 0x03,
		0x00, 0x0a, 'g', 'o', 'o', 'g', 'l', 'e', '.', 'c', 'o', 'm',
		0x00, 0x00, 0x01, 0x11,

		0x00, 0x00, 0x00, 0x00}

	topicsNoBrokersMetadataResponse = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x02,

		0x00, 0x00,
		0x00, 0x03, 'f', 'o', 'o',
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x04,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x07,
		0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00,

		0x00, 0x00,
		0x00, 0x03, 'b', 'a', 'r',
		0x00, 0x00, 0x00, 0x00}
)

func TestEmptyMetadataResponse(t *testing.T) {
	response := MetadataResponse{}

	testDecodable(t, "empty", &response, emptyMetadataResponse)
	if len(response.Brokers) != 0 {
		t.Error("Decoding produced", len(response.Brokers), "brokers where there were none!")
	}
	if len(response.Topics) != 0 {
		t.Error("Decoding produced", len(response.Topics), "topics where there were none!")
	}
}

func TestMetadataResponseWithBrokers(t *testing.T) {
	response := MetadataResponse{}

	testDecodable(t, "brokers, no topics", &response, brokersNoTopicsMetadataResponse)
	if len(response.Brokers) == 2 {
		if response.Brokers[0].id != 0xabff {
			t.Error("Decoding produced invalid broker 0 id.")
		}
		if response.Brokers[0].host != "localhost" {
			t.Error("Decoding produced invalid broker 0 host.")
		}
		if response.Brokers[0].port != 0x33 {
			t.Error("Decoding produced invalid broker 0 port.")
		}
		if response.Brokers[1].id != 0x010203 {
			t.Error("Decoding produced invalid broker 1 id.")
		}
		if response.Brokers[1].host != "google.com" {
			t.Error("Decoding produced invalid broker 1 host.")
		}
		if response.Brokers[1].port != 0x111 {
			t.Error("Decoding produced invalid broker 1 port.")
		}
	} else {
		t.Error("Decoding produced", len(response.Brokers), "brokers where there were two!")
	}
	if len(response.Topics) != 0 {
		t.Error("Decoding produced", len(response.Topics), "topics where there were none!")
	}
}

func TestMetadataResponseWithTopics(t *testing.T) {
	response := MetadataResponse{}

	testDecodable(t, "topics, no brokers", &response, topicsNoBrokersMetadataResponse)
	if len(response.Brokers) != 0 {
		t.Error("Decoding produced", len(response.Brokers), "brokers where there were none!")
	}
	if len(response.Topics) == 2 {
		if response.Topics[0].Err != NO_ERROR {
			t.Error("Decoding produced invalid topic 0 error.")
		}
		if response.Topics[0].Name != "foo" {
			t.Error("Decoding produced invalid topic 0 name.")
		}
		if len(response.Topics[0].Partitions) == 1 {
			if response.Topics[0].Partitions[0].Err != INVALID_MESSAGE_SIZE {
				t.Error("Decoding produced invalid topic 0 partition 0 error.")
			}
			if response.Topics[0].Partitions[0].Id != 0x01 {
				t.Error("Decoding produced invalid topic 0 partition 0 id.")
			}
			if response.Topics[0].Partitions[0].Leader != 0x07 {
				t.Error("Decoding produced invalid topic 0 partition 0 leader.")
			}
			if len(response.Topics[0].Partitions[0].Replicas) == 3 {
				for i := 0; i < 3; i++ {
					if response.Topics[0].Partitions[0].Replicas[i] != int32(i+1) {
						t.Error("Decoding produced invalid topic 0 partition 0 replica", i)
					}
				}
			} else {
				t.Error("Decoding produced invalid topic 0 partition 0 replicas.")
			}
			if len(response.Topics[0].Partitions[0].Isr) != 0 {
				t.Error("Decoding produced invalid topic 0 partition 0 isr length.")
			}
		} else {
			t.Error("Decoding produced invalid partition count for topic 0.")
		}
		if response.Topics[1].Err != NO_ERROR {
			t.Error("Decoding produced invalid topic 1 error.")
		}
		if response.Topics[1].Name != "bar" {
			t.Error("Decoding produced invalid topic 0 name.")
		}
		if len(response.Topics[1].Partitions) != 0 {
			t.Error("Decoding produced invalid partition count for topic 1.")
		}
	} else {
		t.Error("Decoding produced", len(response.Topics), "topics where there were two!")
	}
}
