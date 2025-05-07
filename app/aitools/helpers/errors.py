from typing import Any
from helpers.logger import service_logger

class VapusAiErrorsTemplates:
    """
    A class that provides error templates for Vapus AI.

    This class contains methods to log and handle different types of errors in Vapus AI.
    """
    @staticmethod
    def log_error(e: Any, err_text: str):
        """
        Logs an error message along with the error type and details.

        Args:
            e (Any): The exception object.
            err_text (str): The additional error details.

        Returns:
            None
        """
        service_logger.error("""error: {err_msg} type: {err_type} | details: {err_text}""",
                              errType=type(e).__name__, err_msg=e.__str__(),
                                err_text=err_text)
    @staticmethod
    def prompt_input_error(e: Any):
        """
        Logs the error and generates an error message for input prompt.

        Args:
            e (Any): The error that occurred.

        Returns:
            str: The error message formatted with the error.
        """
        VapusAiErrorsTemplates.log_error(e, em)
        em = "error while generating input prompt for llm"
        return em.format(error=e)
    @staticmethod
    def prompt_output_parse_error(e: Any):
        """
        Handles the error that occurs while parsing the output from LLM response.

        Args:
            e (Any): The error that occurred.

        Returns:
            str: The error message.
        """
        em = "error while parsing output from LLM response"
        VapusAiErrorsTemplates.log_error(e, em)
        return em.format(error=e)
    @staticmethod
    def invalid_llm_config_error(e: Any):
        """
        Raises an error for invalid LLM config provided in the request.

        Args:
            e (Any): The error object.

        Returns:
            str: The error message.
        """
        em = "invalid LLM config provided in the request"
        VapusAiErrorsTemplates.log_error(e, em)
        return em.format(error=e)
    @staticmethod
    def llm_client_error(e: Any):
        """
        Handles the error raised while interacting with the LLM client.

        Args:
            e (Any): The error object raised.

        Returns:
            str: The error message formatted with the error object.
        """
        em = "error while handling llm client"
        VapusAiErrorsTemplates.log_error(e, em)
        return em.format(error=e)
    
class VapusAiError(Exception):
    """
    Custom exception class for Vapus AI errors.

    Attributes:
        message (str): The error message associated with the exception.
    """

    def __init__(self, m):
        self.message = m

    def __str__(self):
        return self.message
        