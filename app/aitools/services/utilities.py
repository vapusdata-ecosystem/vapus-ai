from utils.importer import proto_importer
from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import RecognizerResult, OperatorConfig
from google.protobuf.json_format import MessageToDict

from rake_nltk import Rake
from keybert import KeyBERT
import spacy
import pytextrank
import nltk
import gensim
import re
proto_importer()

import grpc
from protos.vapus_aiutilities.v1alpha1 import vapus_aiutilities_pb2_grpc as aiutilities
import protos.vapus_aiutilities.v1alpha1.vapus_aiutilities_pb2 as pb2
from sentence_transformers import SentenceTransformer
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lex_rank import LexRankSummarizer
from sumy.summarizers.lsa import LsaSummarizer
class Utilities():

    embeddingModel = SentenceTransformer("all-MiniLM-L6-v2")
    analyzer = AnalyzerEngine()
    nltk.download('stopwords')
    nltk.download('punkt')
    nltk.download('punkt_tab')

    def getAnalyzedOutputs(self,result):

        analyzedOutputs = []
        for res in result:
            item_dict = res.to_dict()                   
            analyzedOutput = pb2.AnalyzedOutput()
            analyzedOutput.type = item_dict.get("entity_type")
            analyzedOutput.start = item_dict.get("start")
            analyzedOutput.end  = item_dict.get("end")
            analyzedOutput.score = item_dict.get("score")    
            analyzedOutputs.append(analyzedOutput)
        return analyzedOutputs    

    def GenerateEmbedding(self, request, context):
        
        '''
            logic implementation in services
        '''
        try:
            sentences = request.text 
       
            embeddings = self.embeddingModel.encode(sentences)

            response = pb2.GenerateEmbeddingResponse()
           
            for idx, embedding in enumerate(embeddings):
                embedding_proto = response.Embeddings()
                embedding_proto.embedding.extend(embedding.tolist())  
                embedding_proto.index = idx
                response.embeddings.append(embedding_proto)

            return response

        except Exception as e:
          
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"An error occurred: {str(e)}")
            return pb2.GenerateEmbeddingResponse()
        
    def SensitivityAnalyzer(self, request, context):

        response = pb2.SensitivityAnalyzerResponse()
        
        if request.action == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Invalid action specified")
            return response

        
        if request.action == 1:
            '''
                analyze
            '''
            
            for index,text in enumerate(request.text):
               
                try:
                    result = self.analyzer.analyze(text = text ,language="en",entities=request.entities)
                except Exception as e:
                    raise e
                
                
                processedOutput = response.ProcessedOutput()
                processedOutput.text = text
                processedOutput.index = index
                for entity in request.entities:
                    processedOutput.entities.append(entity)
                processedOutput.action  = request.postDetectAction
                

                # for res in result:

                #     item_dict = res.to_dict()                   
                #     analyzedOutput = pb2.AnalyzedOutput()
                #     analyzedOutput.type = item_dict.get("entity_type")
                #     analyzedOutput.start = item_dict.get("start")
                #     analyzedOutput.end  = item_dict.get("end")
                #     analyzedOutput.score = item_dict.get("score")
                #     processedOutput.AnalyzedOutputs.append(analyzedOutput)
                
                analyzedOutputs = self.getAnalyzedOutputs(result)
                for item in analyzedOutputs:
                    processedOutput.AnalyzedOutputs.append(item)

                response.output.append(processedOutput)
            return response
        if request.action == 2:
            '''
                act
            '''

            if request.postDetectAction == 0:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                context.set_details("Invalid action specified")
                return response

            for index,text in enumerate(request.text):
                result = self.analyzer.analyze(text = text ,language="en",entities=request.entities)
                processedOutput = response.ProcessedOutput()
                processedOutput.text = text
                processedOutput.index = index
                for entity in request.entities:
                    processedOutput.entities.append(entity)
                processedOutput.action  = request.postDetectAction
                analyzer_results = []

                
                # for res in result:

                #     item_dict = res.to_dict()                   
                #     analyzedOutput = pb2.AnalyzedOutput()
                #     analyzedOutput.type = item_dict.get("entity_type")
                #     analyzedOutput.start = item_dict.get("start")
                #     analyzedOutput.end  = item_dict.get("end")
                #     analyzedOutput.score = item_dict.get("score")
                    
                    
                #     if item_dict.get("score")>=0.7:
                #         analyzer_results.append(res)
                    
                #     processedOutput.AnalyzedOutputs.append(analyzedOutput)
                
                analyzedOutputs = self.getAnalyzedOutputs(result)
                for item in analyzedOutputs:
                    processedOutput.AnalyzedOutputs.append(item)
                
                for index,item in enumerate(analyzedOutputs):
                    if item.score >= 0.7:
                        analyzer_results.append(result[index])
                
                operators = {}
                
                placeholder = {1:"xxxx",2:"Placeholder",3:""}

               
                for item in analyzer_results:
                    operators[item.entity_type] = OperatorConfig("replace", {"new_value": placeholder.get(request.postDetectAction)})
                
                engine = AnonymizerEngine()
                editedText = engine.anonymize(text = text,analyzer_results=analyzer_results,operators=operators)
                processedOutput.text = editedText.text
                response.output.append(processedOutput)
            return response
        
    # def SummarizeText(self, request, context):
    #     response = pb2.SummarizerResponse()
    #     text = request.text.decode('utf-8')
        
       
    #     rake = Rake()
    #     rake.extract_keywords_from_text(text)
    #     rake_keywords = set(rake.get_ranked_phrases()[:10]) 
        
        
    #     kw_model = KeyBERT()
    #     bert_results = kw_model.extract_keywords(text, keyphrase_ngram_range=(1, 2), stop_words='english', top_n=10)
    #     bert_keywords = set([kw for kw, score in bert_results])
        
        
    #     nlp = spacy.load("en_core_web_sm")
        
    #     if "textrank" not in nlp.pipe_names:
    #         nlp.add_pipe("textrank")
    #     doc = nlp(text)
    #     textrank_keywords = set([phrase.text for phrase in doc._.phrases[:10]])
        
        
    #     combined_keywords = list(rake_keywords.union(bert_keywords).union(textrank_keywords))
        
        
    #     keyword_scores = {}
    #     for keyword in combined_keywords:
    #         score = 0
    #         if keyword in rake_keywords:
    #             score += 1
    #         if keyword in bert_keywords:
    #             score += 1
    #         if keyword in textrank_keywords:
    #             score += 1
    #         keyword_scores[keyword] = score
        
        
    #     sorted_keywords = sorted(keyword_scores.items(), key=lambda x: (-x[1], x[0]))
    #     top_keywords = [kw for kw, score in sorted_keywords][:10]  # final top 10 keywords

        
    #     response.data.extend([kw.encode('utf-8') for kw in top_keywords])
        
    #     return response

    # def SummarizeText(self, request, context):
    #     response = pb2.SummarizerResponse()
    #     text = request.text.decode('utf-8')
        
    #     try:
    #         # The ratio parameter defines the fraction of sentences to keep.
    #         summary = summarize(text, ratio=0.1)  
    #     except Exception as e:
    #         # In case summarization fails (e.g., text is too short), fall back to original text.
    #         summary = text

    #     response.data.append(summary.encode('utf-8'))
    #     return response


    # def SummarizeText(self, request, context):
    #     response = pb2.SummarizerResponse()
    #     text = request.text.decode('utf-8')
        
       
    #     cleaned_text = re.sub(r"[^A-Za-z0-9\s\.\,\?\!\:\;\-\(\)\'\"]", "", text)
        
    #     try:

    #         summary = summarize(cleaned_text, ratio=0.2)
    #     except Exception as e:

    #         summary = cleaned_text

    #     response.data.append(summary.encode('utf-8'))
    #     return response

    # def SummarizeText(self, request, context):
    #     response = pb2.SummarizerResponse()
    #     text = request.text.decode('utf-8')
        
      
    #     cleaned_text = re.sub(r"[^A-Za-z0-9\s\.\,\?\!\:\;\-\(\)\'\"]", "", text)
        
       
    #     parser = PlaintextParser.from_string(cleaned_text, Tokenizer("english"))
        
        
    #     summarizer = LexRankSummarizer()
        
        
    #     summary_sentences = summarizer(parser.document, sentences_count=4)
        
        
    #     summary = " ".join([str(sentence) for sentence in summary_sentences])
        
    #     response.data.append(summary.encode('utf-8'))
    #     return response
    
    def SummarizeText(self, request, context):
        print("SummarizeText called",request)
        response = pb2.SummarizerResponse()
        text = request.text.decode('utf-8')
        
        
        text = re.sub(r"(Input\s*::\s*\[.*?\])|(Output\s*::\s*\[.*?\])", "", text, flags=re.IGNORECASE)
        
        
        cleaned_text = re.sub(r"[^A-Za-z0-9\s\.\,\?\!\:\;\-\(\)\'\"]", "", text)
        
        
        parser = PlaintextParser.from_string(cleaned_text, Tokenizer("english"))
        
       
        total_sentences = len(list(parser.document.sentences))
        sentences_count = request.sentences if total_sentences >= 8 else max(1, total_sentences // 2)
      
        lexrank = LexRankSummarizer()
        lsa = LsaSummarizer()
        
        lexrank_summary = [str(sentence) for sentence in lexrank(parser.document, sentences_count)]
        lsa_summary = [str(sentence) for sentence in lsa(parser.document, sentences_count)]
        
        combined = []
        for sentence in lexrank_summary + lsa_summary:
            if sentence not in combined:
                combined.append(sentence)
        final_summary = " ".join(combined[:sentences_count])
        response.data.append(final_summary.encode('utf-8'))
        return response








        








