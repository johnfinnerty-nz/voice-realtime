"""Tests for the minimal OpenAI Realtime client."""

from __future__ import annotations

import base64
import json
import tempfile
import unittest
from pathlib import Path

from examples.python_client import realtime_client


class MockTransport:
    def __init__(self, events: list[dict[str, object]]) -> None:
        self.events = iter(json.dumps(event) for event in events)

    def recv(self) -> str:
        return next(self.events)


class ReceiveResponseTest(unittest.TestCase):
    def test_saves_audio_after_successful_response(self) -> None:
        for response in ({"status": "completed"}, {}):
            with self.subTest(response=response), tempfile.TemporaryDirectory() as directory:
                output_path = Path(directory) / "response.pcm"
                transport = MockTransport(
                    [
                        {
                            "type": "response.audio.delta",
                            "delta": base64.b64encode(b"complete audio").decode("ascii"),
                        },
                        {"type": "response.done", "response": response},
                    ]
                )

                realtime_client.receive_response(  # type: ignore[arg-type]
                    transport,
                    output_path,
                )

                self.assertEqual(output_path.read_bytes(), b"complete audio")

    def test_does_not_save_partial_audio_for_unsuccessful_response(self) -> None:
        cases = {
            "cancelled": {"type": "cancelled", "reason": "turn_detected"},
            "failed": {
                "type": "failed",
                "error": {
                    "type": "server_error",
                    "code": "provider_failed",
                    "message": "provider rejected request",
                },
            },
            "incomplete": {"type": "incomplete", "reason": "max_output_tokens"},
        }

        for status, status_details in cases.items():
            with self.subTest(status=status), tempfile.TemporaryDirectory() as directory:
                output_path = Path(directory) / "response.pcm"
                transport = MockTransport(
                    [
                        {
                            "type": "response.audio.delta",
                            "delta": base64.b64encode(b"partial audio").decode("ascii"),
                        },
                        {
                            "type": "response.done",
                            "response": {
                                "status": status,
                                "status_details": status_details,
                            },
                        },
                    ]
                )

                with self.assertRaisesRegex(
                    RuntimeError,
                    f"Response ended with status {status}",
                ) as error:
                    realtime_client.receive_response(  # type: ignore[arg-type]
                        transport,
                        output_path,
                    )

                self.assertIn(json.dumps(status_details), str(error.exception))
                self.assertFalse(output_path.exists())


if __name__ == "__main__":
    unittest.main()
