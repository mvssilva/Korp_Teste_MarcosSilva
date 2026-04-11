import { Component } from '@angular/core';
import { CommonModule } from '@angular/common'; 
import { FormsModule } from '@angular/forms'; 

@Component({
  selector: 'app-estoque',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './estoque.component.html',
  styleUrl: './estoque.component.css'
})
export class EstoqueComponent {
  novoCodigo: string = '';
  novaDescricao: string = '';
  novoSaldo: number = 0;

  listaProdutos = [
    { codigo: 'P001', descricao: 'Teclado Mecânico', saldo: 13 }
  ];

  salvar() {
    this.listaProdutos.push({
      codigo: this.novoCodigo,
      descricao: this.novaDescricao,
      saldo: this.novoSaldo
    });
    this.novoCodigo = '';
    this.novaDescricao = '';
    this.novoSaldo = 0;
  }
}