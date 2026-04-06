package main

import (
	"flag"
	"fmt"
	"strings"
	
)
import json "encoding/json"
import os "os"


type Tarefa struct {
	ID     int
	Nome   string
	Status string
}

const (
	paraFazer = "Para fazer"
	andamento = "Em andamento"
	concluida = "Concluída"
)

func main() {

	tarefas := carregarTarefas()

	
	lista := flag.Int("l", 0, "Listar por status (1: Para fazer, 2: Em andamento, 3: Concluída)")
	listarTodas := flag.Bool("all", false, "Listar todas as tarefas")
	add := flag.String("a", "", "Adicionar tarefa (ex: 'Comprar pão:1')")
	move := flag.String("m", "", "Mover tarefa (ex: '1:2')")
	remove := flag.Int("r", 0, "Remover tarefa por ID")
	
	flag.Parse()

	
	if flag.NFlag() == 0 {
		fmt.Println("Bem-vindo ao To-Do. Use -help para ver os comandos.")
		return
	}

	
	if *add != "" {
		
		parts := strings.Split(*add, ":") 
		if len(parts) != 2 {
			fmt.Println("Erro: Formato inválido. Use 'nome:statusID' (ex: 'Pão:1').")
			return
		}
		
		nome := parts[0]
		statusStr := parts[1]
		statusFormatado := converterStatus(statusStr)

		if statusFormatado == "" {
			fmt.Println("Erro: Status inválido. Use 1, 2 ou 3.")
			return
		}

		tarefas = adicionarTarefa(tarefas, nome, statusFormatado)
		fmt.Printf("Sucesso: Tarefa '%s' criada em '%s'.\n", nome, statusFormatado)
	}

	if *move != "" {
		parts := strings.Split(*move, ":")
		if len(parts) != 2 {
			fmt.Println("Erro: Formato inválido. Use 'idDaTarefa:novoStatusID' (ex: '1:2').")
			return
		}

		var id int
		fmt.Sscanf(parts[0], "%d", &id)
		novoStatus := converterStatus(parts[1])

		if novoStatus != "" {
			tarefas = moverTarefa(tarefas, id, novoStatus)
			fmt.Printf("Sucesso: Tarefa %d movida para '%s'.\n", id, novoStatus)
		}
	}


	if *listarTodas {
		fmt.Println("=== TODAS AS TAREFAS ===")
		exibirTarefas(tarefas, "") 
	} else if *lista != 0 {
		statusStr := fmt.Sprintf("%d", *lista) 
		status := converterStatus(statusStr)
		if status != "" {
			fmt.Printf("=== TAREFAS: %s ===\n", strings.ToUpper(status))
			exibirTarefas(tarefas, status)
		}
	}
	if *remove != 0 {
		tarefas = removerTarefa(tarefas, *remove)
		fmt.Printf("Sucesso: Tarefa %d removida.\n", *remove)
	}
	salvarTarefas(tarefas)
}

func converterStatus(codigo string) string {
	switch codigo {
	case "1":
		return paraFazer
	case "2":
		return andamento
	case "3":
		return concluida
	default:
		return ""
	}
}

func adicionarTarefa(tarefas []Tarefa, nome string, status string) []Tarefa {
	id := len(tarefas) + 1
	tarefa := Tarefa{ID: id, Nome: nome, Status: status}
	return append(tarefas, tarefa)
}

func moverTarefa(tarefas []Tarefa, id int, novoStatus string) []Tarefa {
	for i, tarefa := range tarefas {
		if tarefa.ID == id {
			tarefas[i].Status = novoStatus
			break 
		}
	}
	return tarefas
}

func exibirTarefas(tarefas []Tarefa, statusFiltro string) {
	for _, tarefa := range tarefas {
		if statusFiltro == "" || tarefa.Status == statusFiltro {
			fmt.Printf("[%d] %s -> %s\n", tarefa.ID, tarefa.Nome, tarefa.Status)
		}
	}
}

func removerTarefa(tarefas []Tarefa, id int) []Tarefa {
	for i, tarefa := range tarefas {
		if tarefa.ID == id {
			return append(tarefas[:i], tarefas[i+1:]...)
		}
	}
	return tarefas
}

func salvarTarefas(tarefas []Tarefa) {
	jsonData, err := json.MarshalIndent(tarefas, "", "  ")
	if err != nil {
		fmt.Println("Erro ao salvar tarefas:", err)
		return
	}
	
	err = os.WriteFile("tarefas.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Erro ao escrever arquivo:", err)
	}
}

func carregarTarefas() []Tarefa	{
	var tarefas []Tarefa
	jsonData, err := os.ReadFile("tarefas.json")
	if err != nil {
		if os.IsNotExist(err) {
			return tarefas 
		}
		fmt.Println("Erro ao ler arquivo:", err)
		return tarefas
	}

	err = json.Unmarshal(jsonData, &tarefas)
	if err != nil {
		fmt.Println("Erro ao carregar tarefas:", err)
	}
	return tarefas
}