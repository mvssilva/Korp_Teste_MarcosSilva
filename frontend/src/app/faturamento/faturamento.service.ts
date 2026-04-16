import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class FaturamentoService {
  
  private apiUrlFaturamento = 'http://localhost:8081/notas';
  private apiUrlEstoque = 'http://localhost:8080/produtos'; 

  constructor(private http: HttpClient) { }

  // Cria a nota com status Aberta
  criarNota(dadosFaturamento: any): Observable<any> {
    return this.http.post(this.apiUrlFaturamento, dadosFaturamento);
  }

  // Consulta se o produto existe no estoque e traz o nome dele
  consultarProdutoNoEstoque(codigo: string): Observable<any> {
    return this.http.get(`${this.apiUrlEstoque}/${codigo}`);
  }

  // Avisa o Go para imprimir e dar baixa no estoque
  imprimirNota(idNota: string | number): Observable<any> {
    return this.http.post(`${this.apiUrlFaturamento}/${idNota}/imprimir`, {});
  }
}